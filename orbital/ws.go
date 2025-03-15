package orbital

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coder/websocket"
	"net/http"
	"orbital/pkg/logger"
	"orbital/pkg/proto"
	"time"
)

type (
	HandlerFunc func(connID string, data []byte)

	Topic struct {
		Name    string
		Handler HandlerFunc
	}

	WsMetadata struct {
		Topic string `json:"topic"`
	}

	WsService interface {
		Register(topic Topic)
		Broadcast(m proto.Message)
		SendTo(connectionID string, m proto.Message) error
	}

	WsConn struct {
		log               *logger.Logger
		topics            map[string]Topic
		connectionManager *WsConnectionManager
	}
)

func (ws *WsConn) Register(topic Topic) {
	ws.log.Info("Register topic", "topic", topic.Name)

	if _, found := ws.topics[topic.Name]; found {
		ws.log.Error(fmt.Sprintf("topic %s is already registered", topic.Name))
		return
	}

	ws.topics[topic.Name] = topic
}

func (ws *WsConn) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

	ws.handleConnection(wsConn)
}

func (ws *WsConn) Broadcast(m proto.Message) {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	raw, err := json.Marshal(m)
	if err != nil {
		fmt.Println("cannot marshal message for broadcast")
		return
	}

	ws.connectionManager.Broadcast(ctx, raw)
}

func (ws *WsConn) SendTo(connectionID string, m proto.Message) error {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	raw, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return ws.connectionManager.SendTo(ctx, connectionID, raw)

}

func (ws *WsConn) handleConnection(conn *websocket.Conn) {
	connID := genConnID()
	ws.connectionManager.AddConnection(connID, conn)

	defer func() {
		fmt.Println("close connection")
		ws.connectionManager.RemoveConnection(connID)
		conn.Close(websocket.StatusNormalClosure, "closing connection")
	}()

	ws.log.Info("New connection joined", "connID", connID)

	for {
		_, msg, err := conn.Read(context.Background())
		if err != nil {
			ws.log.Error(err.Error(), "connection", "read error", "resolution", "closing connection")
			return
		}

		var message proto.Message
		if err = json.Unmarshal(msg, &message); err != nil {
			ws.log.Error(err.Error(), "connection", "unmarshal error", "resolution", "skip message")
			continue
		}

		var metadata *WsMetadata
		if err = json.Unmarshal(message.Metadata, &metadata); err != nil {
			ws.log.Warn("metadata cannot be decoded")
			continue
		}

		handler, found := ws.topics[metadata.Topic]
		if !found {
			ws.log.Error("topic field is missing", "topic", metadata.Topic)
			continue
		}

		handler.Handler(connID, msg)
	}
}

func NewWsConn(log *logger.Logger) *WsConn {
	wsConn := &WsConn{
		log:               log,
		topics:            make(map[string]Topic),
		connectionManager: NewWsConnectionManager(),
	}

	return wsConn
}
