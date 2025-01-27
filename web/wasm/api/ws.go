package api

import (
	"encoding/json"
	"fmt"
	"orbital/pkg/proto"
	"orbital/web/wasm/dom"
	"syscall/js"
)

type (
	WsMetadata struct {
		Topic string `json:"topic"`
	}
	HandlerFunc func(data []byte)
)

type WsConn struct {
	client       js.Value
	topics       map[string]HandlerFunc
	isOpen       bool
	allowsBinary bool
}

func (ws *WsConn) IsOpen() bool {
	return ws.isOpen
}

func (ws *WsConn) Send(msg proto.Message) {
	if !ws.isOpen {
		dom.PrintToConsole("WebSocket closed")
		return
	}

	raw, err := json.Marshal(msg)
	if err != nil {
		dom.PrintToConsole(err.Error())
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
		dom.PrintToConsole("WebSocket codec switch to binary")
		socket.Set("binaryType", "arraybuffer")
	}

	ws.client.Call("addEventListener", "open", js.FuncOf(ws.onOpen))
	ws.client.Call("addEventListener", "close", js.FuncOf(ws.onClose))
	ws.client.Call("addEventListener", "message", js.FuncOf(ws.onMessage))
	ws.client.Call("addEventListener", "error", js.FuncOf(ws.onError))
}

func (ws *WsConn) onOpen(this js.Value, args []js.Value) interface{} {
	ws.isOpen = true
	dom.PrintToConsole("WebSocket connection open")
	return nil
}

func (ws *WsConn) onClose(this js.Value, args []js.Value) interface{} {
	ws.isOpen = false
	dom.PrintToConsole("WebSocket connection closed")
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
	dom.PrintToConsole("WebSocket connection error")
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
		dom.PrintToConsole("[routeMessage] not valid JSON, fallback to raw textData")
		return
	}

	ok, err := msg.Verify()
	if err != nil {
		dom.PrintToConsole(err.Error())
		return
	}

	if !ok {
		dom.PrintToConsole("[routeMessage] message signature not valid")
		return
	}

	var metadata WsMetadata
	if err = json.Unmarshal(msg.Metadata, &metadata); err != nil {
		dom.PrintToConsole("[routeMessage] not valid msg.Metadata format, fallback to raw textData")
		return
	}

	handler, exists := ws.topics[metadata.Topic]
	if !exists {
		dom.PrintToConsole("[routeMessage] topic not found", metadata.Topic)
		return
	}

	handler(msg.Body)
}

func NewWsConn(binaryMode bool) *WsConn {

	wsConn := &WsConn{
		topics:       make(map[string]HandlerFunc),
		isOpen:       false,
		allowsBinary: binaryMode,
	}
	wsConn.init()

	return wsConn
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
