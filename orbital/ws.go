package orbital

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"orbital/pkg/cryptographer"
	"orbital/pkg/logger"
	"time"

	"github.com/coder/websocket"
)

type (
	HandlerFunc func(ctx context.Context, connID string, data []byte)

	Topic struct {
		Name    string
		Handler HandlerFunc
	}

	WsService interface {
		SetSecretKey(secretKey cryptographer.PrivateKey)
		Register(topic Topic)
		Broadcast(ctx context.Context, m cryptographer.Message)
		SendTo(ctx context.Context, connectionID string, m cryptographer.Message) error
		ServeHTTP(w http.ResponseWriter, r *http.Request)
	}

	WsConn struct {
		secretKey         cryptographer.PrivateKey
		log               *logger.Logger
		topics            map[string]Topic
		connectionManager *WsConnectionManager
		writeTimeout      time.Duration
		idleTimeout       time.Duration
	}
)

func (ws *WsConn) SetSecretKey(secretKey cryptographer.PrivateKey) {
	ws.secretKey = secretKey
}

func (ws *WsConn) Register(topic Topic) {
	ws.log.Info("Register topic", "topic", topic.Name)

	if _, found := ws.topics[topic.Name]; found {
		ws.log.Error(fmt.Sprintf("topic %s is already registered", topic.Name))
		return
	}

	ws.topics[topic.Name] = topic
}

func (ws *WsConn) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var (
		wsConn *websocket.Conn
		err    error
	)

	wsConn, err = websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		_ = Encode(w, r, http.StatusInternalServerError, Error{
			Internal,
			err,
		})
		return
	}

	ws.handleConnection(ctx, wsConn)
}

func (ws *WsConn) Broadcast(ctx context.Context, m cryptographer.Message) {

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, ws.writeTimeout)
		defer cancel()
	}

	raw, err := json.Marshal(m)
	if err != nil {
		fmt.Println("cannot marshal message for broadcast")
		return
	}

	ws.connectionManager.Broadcast(ctx, raw)
}

func (ws *WsConn) SendTo(ctx context.Context, connID string, m cryptographer.Message) error {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	raw, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return ws.connectionManager.SendTo(ctx, connID, raw)

}

func (ws *WsConn) handleConnection(ctx context.Context, conn *websocket.Conn) {
	connID := genConnID()
	ws.connectionManager.AddConnection(connID, conn)

	defer func() {
		ws.connectionManager.RemoveConnection(connID)
		_ = conn.Close(websocket.StatusNormalClosure, "closing connection")
	}()

	ws.log.Info("New connection joined", "connID", connID)

	connCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Adjust the pingInterval in case idleTimeout is set
	pingInterval := 30 * time.Second
	if ws.idleTimeout > 0 && ws.idleTimeout/2 < pingInterval {
		pingInterval = ws.idleTimeout / 2
	}

	// Start heartbeat
	go keepAlive(connCtx, conn, pingInterval, 5*time.Second)

	// Welcome the client
	ws.sendWelcomeMessage(connCtx, connID)

	for {
		readCtx := connCtx
		var cancelRead context.CancelFunc
		if ws.idleTimeout > 0 {
			readCtx, cancelRead = context.WithTimeout(context.Background(), ws.idleTimeout)
		}
		_, msg, err := conn.Read(readCtx)
		if cancelRead != nil {
			cancelRead()
		}

		if err != nil {
			ws.log.Error(err.Error(), "connection", "read error", "resolution", "closing connection")
			return
		}

		var message cryptographer.Message
		if err = json.Unmarshal(msg, &message); err != nil {
			ws.log.Error(err.Error(), "connection", "unmarshal error", "resolution", "skip message")
			continue
		}

		// TODO: Add verify message

		t, err := topic(message.Metadata.Domain, message.Metadata.Action)
		if err != nil {
			ws.log.Error(err.Error(), "topic parsing error")
			continue
		}

		handler, found := ws.topics[t]
		if !found {
			ws.log.Error("topic field is missing", "topic", t)
			continue
		}

		func() {
			defer func() {
				if r := recover(); r != nil {
					ws.log.Error("recovered from panic", "topic", t, "err", r)
				}
			}()
			handler.Handler(connCtx, connID, msg)
		}()
	}
}

type WelcomeMessage struct {
	ConnID     string         `json:"connId"`
	ServerTime int64          `json:"serverTime"`
	Code       Code           `json:"code"`
	Error      *ErrorResponse `json:"error,omitempty"`
}

func (ws *WsConn) sendWelcomeMessage(ctx context.Context, connID string) {
	msg, err := cryptographer.Encode(ws.secretKey, cryptographer.Metadata{
		Domain: "system",
		Action: "welcome",
	}, WelcomeMessage{
		Code:       OK,
		ConnID:     connID,
		ServerTime: time.Now().Unix(),
	})
	if err != nil {
		ws.log.Error(err.Error(), "welcome message encoding error")
		return
	}

	if err = ws.SendTo(ctx, connID, *msg); err != nil {
		ws.log.Error(err.Error(), "welcome message sending error")
		return
	}
}

func NewWsConn(log *logger.Logger) *WsConn {
	wsConn := &WsConn{
		log:               log,
		topics:            make(map[string]Topic),
		connectionManager: NewWsConnectionManager(),
		writeTimeout:      5 * time.Second,
		idleTimeout:       30 * time.Second,
	}

	return wsConn
}

func topic(domain, action string, cid ...string) (string, error) {
	if domain == "" {
		return "", errors.New("domain is required")
	}

	if action == "" {
		return "", errors.New("action is required")
	}

	t := fmt.Sprintf("%s/%s", domain, action)
	if len(cid) != 0 && cid[0] != "" {
		t = fmt.Sprintf("%s/%s", t, cid[0])
	}

	return t, nil
}

func keepAlive(ctx context.Context, conn *websocket.Conn, interval, timeout time.Duration) {
	tick := time.NewTicker(interval)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			kaCtx, cancel := context.WithTimeout(ctx, timeout)
			if err := conn.Ping(kaCtx); err != nil {
				cancel()

				_ = conn.Close(websocket.StatusPolicyViolation, "keep alive error")
				return
			}
			cancel()
		}
	}
}
