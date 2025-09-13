package orbital

import (
	"context"
	"fmt"
	"orbital/pkg/stringer"
	"strings"
	"sync"

	"github.com/coder/websocket"
)

type WsConnection struct {
	ID     string
	Conn   *websocket.Conn
	UserID string // Custom set by the user
}

type WsConnectionManager struct {
	mu          sync.RWMutex
	connections map[string]*WsConnection
}

func (wcm *WsConnectionManager) AddConnection(id string, conn *websocket.Conn) {
	wcm.mu.Lock()
	defer wcm.mu.Unlock()
	wcm.connections[id] = &WsConnection{
		ID:   id,
		Conn: conn,
	}
}

func (wcm *WsConnectionManager) RemoveConnection(id string) {
	wcm.mu.Lock()
	defer wcm.mu.Unlock()

	delete(wcm.connections, id)
}

func (wcm *WsConnectionManager) GetConnection(id string) (*WsConnection, bool) {
	wcm.mu.RLock()
	defer wcm.mu.RUnlock()
	conn, ok := wcm.connections[id]

	return conn, ok
}

func (wcm *WsConnectionManager) SetUserID(id, userID string) {
	wcm.mu.Lock()
	defer wcm.mu.Unlock()

	if c, exists := wcm.connections[id]; exists {
		c.UserID = userID
	}
}

func (wcm *WsConnectionManager) Broadcast(ctx context.Context, message []byte) {
	wcm.mu.RLock()
	defer wcm.mu.RUnlock()

	for id, conn := range wcm.connections {
		if err := conn.Conn.Write(ctx, websocket.MessageBinary, message); err != nil {
			// TODO: Check if this is a correct approach
			wcm.RemoveConnection(id)
		}
	}
}

func (wcm *WsConnectionManager) SendTo(ctx context.Context, id string, message []byte) error {
	conn, exists := wcm.GetConnection(id)
	if !exists {
		return fmt.Errorf("connection not found for id: %s", id)
	}

	return conn.Conn.Write(ctx, websocket.MessageBinary, message)
}

// NewWsConnectionManager create a new connection manager
func NewWsConnectionManager() *WsConnectionManager {
	return &WsConnectionManager{
		connections: make(map[string]*WsConnection),
	}
}

func genConnID() string {
	randStr, _ := stringer.Random(16, stringer.RandNumber, stringer.RandLowercase)
	return strings.Join([]string{"orb", randStr}, ".")
}
