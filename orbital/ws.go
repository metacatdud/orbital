package orbital

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"nhooyr.io/websocket"
	"orbital/pkg/cryptographer"
	"orbital/pkg/logger"
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
)

type WsService interface {
	Register(topic Topic)
	Broadcast(message []byte)
	SendTo(connectionID string, message []byte) error
}

type WsConn struct {
	log               *logger.Logger
	topics            map[string]Topic
	connectionManager *WsConnectionManager
}

func (ws *WsConn) Register(topic Topic) {
	ws.log.Info("Register topic:", "topic", topic.Name)

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

func (ws *WsConn) Broadcast(message []byte) {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	ws.connectionManager.Broadcast(ctx, message)
}

func (ws *WsConn) SendTo(connectionID string, message []byte) error {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	return ws.connectionManager.SendTo(ctx, connectionID, message)

}

func (ws *WsConn) handleConnection(conn *websocket.Conn) {
	connID := genConnID()
	ws.connectionManager.AddConnection(connID, conn)

	defer func() {
		fmt.Println("close connection")
		ws.connectionManager.RemoveConnection(connID)
		conn.Close(websocket.StatusNormalClosure, "closing connection")
	}()

	ws.log.Info("New connection created", "connID", connID)

	for {
		_, msg, err := conn.Read(context.Background())
		if err != nil {
			ws.log.Error(err.Error(), "connection", "read error", "resolution", "closing connection")
			return
		}

		var message cryptographer.Message
		if err = json.Unmarshal(msg, &message); err != nil {
			ws.log.Error(err.Error(), "connection", "unmarshal error", "resolution", "skip message")
			continue
		}

		var metadata WsMetadata
		if err = json.Unmarshal(message.Metadata, &metadata); err != nil {
			ws.log.Warn("topic field is missing")
			continue
		}

		handler, found := ws.topics[metadata.Topic]
		if !found {
			ws.log.Error("topic field is missing", "topic", metadata.Topic)
			return
		}

		handler.Handler(connID, msg)
	}
}

func NewWsConn() *WsConn {
	lg := logger.New(logger.LevelDebug, logger.FormatString)
	wsConn := &WsConn{
		log:               lg,
		topics:            make(map[string]Topic),
		connectionManager: NewWsConnectionManager(),
	}

	// Default topics. Overwrite with your custom one
	wsConn.Register(Topic{
		Name:    TopicOrbitalAuthentication,
		Handler: wsConn.defaultOrbitalLogin,
	})

	return wsConn
}

// NewTopicMessage helper for creating a message
func NewTopicMessage(topic string, data []byte) *cryptographer.Message {
	t := &WsMetadata{
		Topic: topic,
	}
	tBytes, _ := json.Marshal(t)

	return &cryptographer.Message{
		V:         1,
		Timestamp: cryptographer.Now(),
		Metadata:  tBytes,
		Body:      data,
	}
}
