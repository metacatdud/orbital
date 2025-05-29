package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"orbital/pkg/proto"
	"orbital/web/wasm/pkg/dom"
	"syscall/js"
	"time"
)

type HandlerFunc func(data []byte)

type WsConn struct {
	client       js.Value
	topics       map[string]HandlerFunc
	isOpen       bool
	allowsBinary bool

	reconnect            bool
	reconnectAttempts    int
	maxReconnectAttempts int
	reconnectInterval    time.Duration
	reconnectInProgress  bool
}

func NewWsConn(binaryMode bool) *WsConn {

	wsConn := &WsConn{
		topics:               make(map[string]HandlerFunc),
		isOpen:               false,
		allowsBinary:         binaryMode,
		reconnect:            true,
		reconnectInterval:    5 * time.Second,
		maxReconnectAttempts: 3,
	}

	wsConn.init()

	return wsConn
}

func (ws *WsConn) IsOpen() bool {
	return ws.isOpen
}

func (ws *WsConn) Send(msg proto.Message) {
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

func (ws *WsConn) init() {
	wsURL := createWebSocketURL()
	socket := js.Global().Get("WebSocket").New(wsURL)

	ws.client = socket

	if ws.allowsBinary {
		dom.ConsoleLog("WebSocket codec switch to binary")
		socket.Set("binaryType", "arraybuffer")
	}

	ws.client.Call("addEventListener", "open", js.FuncOf(ws.onOpen))
	ws.client.Call("addEventListener", "close", js.FuncOf(ws.onClose))
	ws.client.Call("addEventListener", "message", js.FuncOf(ws.onMessage))
	ws.client.Call("addEventListener", "error", js.FuncOf(ws.onError))
}

func (ws *WsConn) onOpen(this js.Value, args []js.Value) interface{} {
	ws.isOpen = true
	ws.reconnectAttempts = 0

	dom.ConsoleLog("WebSocket connection open")
	return nil
}

func (ws *WsConn) onClose(this js.Value, args []js.Value) interface{} {
	ws.isOpen = false

	dom.ConsoleLog("WebSocket connection closed")
	if ws.reconnect {
		ws.scheduleReconnect()
	}

	return nil
}

func (ws *WsConn) onMessage(this js.Value, args []js.Value) interface{} {
	event := args[0]
	dataVal := event.Get("data")

	if ws.allowsBinary && dataVal.InstanceOf(js.Global().Get("ArrayBuffer")) {
		ws.handleBinaryMessage(dataVal)
		return nil
	}

	ws.handleTextMessage(dataVal)
	return nil
}

func (ws *WsConn) onError(this js.Value, args []js.Value) interface{} {
	dom.ConsoleLog("WebSocket connection error")
	if ws.reconnect {
		ws.scheduleReconnect()
	}
	return nil
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
	length := uint8Array.Get("length").Int()
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

	var msg *proto.Message
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

	var t string
	if msg.Metadata != nil {
		dom.ConsoleLog(msg.Metadata)
		t, err = topic(msg.Metadata.Domain, msg.Metadata.Action, msg.Metadata.CorrelationID)
		if err != nil {
			dom.ConsoleLog(err.Error())
			return
		}
	}

	handler, exists := ws.topics[t]
	if !exists {
		dom.ConsoleLog("[routeMessage] topic not found", t)
		return
	}

	handler(msg.Body)
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

	// Wait the "delay" and call init to reconnect
	time.AfterFunc(delay, func() {
		ws.init()
		ws.reconnectInProgress = false
	})
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
