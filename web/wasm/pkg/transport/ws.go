package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"orbital/pkg/cryptographer"
	"orbital/web/wasm/pkg/dom"
	"sync"
	"syscall/js"
	"time"
)

type HandlerFunc func(data []byte)

type WsConn struct {
	mu                                          sync.Mutex
	client                                      js.Value
	topics                                      map[string]HandlerFunc
	isOpen                                      bool
	allowsBinary                                bool
	reconnect                                   bool
	reconnectAttempts                           int
	maxReconnectAttempts                        int
	reconnectInterval                           time.Duration
	reconnectInProgress                         bool
	heartbeatInterval                           time.Duration
	heartbeatWait                               time.Duration
	lastPong                                    time.Time
	onOpenFn, onCloseFn, onMessageFn, onErrorFn js.Func
	kaCancelFn                                  context.CancelFunc
}

func NewWsConn(binaryMode bool) *WsConn {

	wsConn := &WsConn{
		topics:               make(map[string]HandlerFunc),
		isOpen:               false,
		allowsBinary:         binaryMode,
		reconnect:            true,
		reconnectInterval:    5 * time.Second,
		maxReconnectAttempts: 10,
		heartbeatInterval:    15 * time.Second,
		heartbeatWait:        (15 * time.Second) * 2,
	}

	wsConn.connect()

	return wsConn
}

func (ws *WsConn) IsOpen() bool {
	return ws.isOpen
}

func (ws *WsConn) Send(msg cryptographer.Message) {
	if !ws.isOpen {
		dom.ConsoleWarn("WebSocket closed")
		return
	}

	raw, err := json.Marshal(msg)
	if err != nil {
		dom.ConsoleError(err.Error())
		return
	}

	if ws.allowsBinary {
		ws.sendBinary(raw)
		return
	}

	ws.sendText(raw)
}

func (ws *WsConn) On(topic string, handler HandlerFunc) {
	ws.topics[topic] = handler
}

func (ws *WsConn) connect() {
	wsURL := createWebSocketURL()
	socket := js.Global().Get("WebSocket").New(wsURL)
	ws.client = socket

	if ws.allowsBinary {
		dom.ConsoleLog("WebSocket codec switch to binary")
		socket.Set("binaryType", "arraybuffer")
	}

	ws.onOpenFn = js.FuncOf(ws.onOpen)
	ws.onCloseFn = js.FuncOf(ws.onClose)
	ws.onMessageFn = js.FuncOf(ws.onMessage)
	ws.onErrorFn = js.FuncOf(ws.onError)

	socket.Call("addEventListener", "open", ws.onOpenFn)
	socket.Call("addEventListener", "close", ws.onCloseFn)
	socket.Call("addEventListener", "message", ws.onMessageFn)
	socket.Call("addEventListener", "error", ws.onErrorFn)
}

func (ws *WsConn) onOpen(_ js.Value, _ []js.Value) any {
	ws.mu.Lock()
	ws.isOpen = true
	ws.reconnectAttempts = 0
	ws.lastPong = time.Now()
	ws.mu.Unlock()

	ws.startKeepAlive()

	return nil
}

func (ws *WsConn) onClose(_ js.Value, _ []js.Value) any {
	ws.stopKeepAlive()

	ws.isOpen = false

	dom.ConsoleWarn("WebSocket connection closed")
	if ws.reconnect {
		ws.scheduleReconnect()
	}

	return nil
}

func (ws *WsConn) onMessage(_ js.Value, args []js.Value) any {
	event := args[0]
	dataVal := event.Get("data")

	if ws.allowsBinary && dataVal.InstanceOf(js.Global().Get("ArrayBuffer")) {
		ws.handleBinaryMessage(dataVal)
		return nil
	}

	ws.handleTextMessage(dataVal)
	return nil
}

func (ws *WsConn) onError(_ js.Value, _ []js.Value) any {
	dom.ConsoleLog("WebSocket connection error")
	if ws.reconnect {
		ws.scheduleReconnect()
	}
	return nil
}

func (ws *WsConn) teardown() {
	if ws.client.Truthy() {
		ws.client.Call("removeEventListener", "open", ws.onOpenFn)
		ws.client.Call("removeEventListener", "close", ws.onCloseFn)
		ws.client.Call("removeEventListener", "message", ws.onMessageFn)
		ws.client.Call("removeEventListener", "error", ws.onErrorFn)
	}

	ws.onOpenFn.Release()
	ws.onCloseFn.Release()
	ws.onMessageFn.Release()
	ws.onErrorFn.Release()
}

func (ws *WsConn) sendBinary(data []byte) {
	uint8Array := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(uint8Array, data)

	ws.client.Call("send", uint8Array.Get("buffer"))
}

func (ws *WsConn) sendText(data []byte) {
	ws.client.Call("send", data)
}

func (ws *WsConn) handleBinaryMessage(dataBuffer js.Value) {
	uint8Array := js.Global().Get("Uint8Array").New(dataBuffer)
	length := uint8Array.Length()
	raw := make([]byte, length)
	js.CopyBytesToGo(raw, uint8Array)

	ws.routeMessage(raw)
}

func (ws *WsConn) handleTextMessage(dataVal js.Value) {
	strData := dataVal.String()
	raw := []byte(strData)

	ws.routeMessage(raw)
}

func (ws *WsConn) routeMessage(raw []byte) {

	var (
		msg cryptographer.Message
		t   string
	)
	if err := json.Unmarshal(raw, &msg); err != nil {
		dom.ConsoleLog("[routeMessage] not valid JSON")
		return
	}

	ok, err := msg.Verify()
	if err != nil {
		dom.ConsoleLog(err.Error())
		return
	}

	if !ok {
		dom.ConsoleLog("[routeMessage] message signature not valid")
		return
	}

	t, err = topic(msg.Metadata.Domain, msg.Metadata.Action, msg.Metadata.CorrelationID)
	if err != nil {
		dom.ConsoleLog(err.Error())
		return
	}

	switch t {
	case "system/welcome":
		// --- move this out and allow app level implementation
		dom.ConsoleLog("[routeMessage] system welcome message", string(msg.Body))
	case "system/keepAlivePing":
		ws.Send(makeKeepAlivePong())
	case "system/keepAlivePong":
		ws.mu.Lock()
		ws.lastPong = time.Now()
		dom.ConsoleLog("[routeMessage] keep alive pong", ws.lastPong)
		ws.mu.Unlock()
	default:
		handler, exists := ws.topics[t]
		if !exists {
			dom.ConsoleLog("[routeMessage] topic not found", t)
			return
		}

		handler(msg.Body)
	}
}

func (ws *WsConn) scheduleReconnect() {
	if ws.reconnectInProgress {
		dom.ConsoleLog("[scheduleReconnect] ws already reconnecting")
		return
	}

	ws.reconnectInProgress = true

	if ws.maxReconnectAttempts != -1 && ws.reconnectAttempts >= ws.maxReconnectAttempts {
		dom.ConsoleLog("[scheduleReconnect] max reconnection attempts reached")
		ws.reconnectInProgress = false
		return
	}

	ws.reconnectAttempts++

	delay := time.Duration(ws.reconnectAttempts) * ws.reconnectInterval
	maxDelay := delay * 2
	jitter := time.Duration(rand.Int63n(int64(maxDelay - delay)))
	delay += jitter

	dom.ConsoleLog("[scheduleReconnect] attempting to reconnect in", delay.String())

	// Wait the "delay" and call connect to reconnect
	time.AfterFunc(delay, func() {
		ws.teardown()
		ws.connect()
		ws.reconnectInProgress = false
	})
}

func (ws *WsConn) startKeepAlive() {
	ws.stopKeepAlive()

	ctx, cancel := context.WithCancel(context.Background())
	ws.kaCancelFn = cancel

	go func() {
		tick := time.NewTicker(ws.heartbeatInterval)
		defer tick.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				m := makeKeepAlivePing()
				ws.Send(m)
			}
		}
	}()
}

func (ws *WsConn) stopKeepAlive() {
	ws.mu.Lock()
	if ws.kaCancelFn != nil {
		ws.kaCancelFn()
		ws.kaCancelFn = nil
	}
	ws.mu.Unlock()
}

// createWebSocketURL determine the URL for websocket
func createWebSocketURL() string {
	location := js.Global().Get("window").Get("location")
	protocol := location.Get("protocol").String()
	hostname := location.Get("hostname").String()
	port := location.Get("port").String()

	wsProtocol := "ws"
	if protocol == "https:" {
		wsProtocol = "wss"
	}

	if port == "" {
		return fmt.Sprintf("%s://%s/ws", wsProtocol, hostname)
	}
	return fmt.Sprintf("%s://%s:%s/ws", wsProtocol, hostname, port)
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

// makeKeepAlivePing creates an unsigned message
// TODO: See what the implications are for this message
func makeKeepAlivePing() cryptographer.Message {
	msg := cryptographer.Message{
		V:         0,
		Timestamp: cryptographer.Now(),
		Metadata: cryptographer.Metadata{
			Domain: "system",
			Action: "keepAlivePing",
		},
		Body: nil,
	}

	return msg
}

func makeKeepAlivePong() cryptographer.Message {
	msg := cryptographer.Message{
		V:         0,
		Timestamp: cryptographer.Now(),
		Metadata: cryptographer.Metadata{
			Domain: "system",
			Action: "keepAlivePong",
		},
		Body: nil,
	}

	return msg
}
