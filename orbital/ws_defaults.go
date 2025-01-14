package orbital

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"orbital/pkg/cryptographer"
)

const (
	TopicOrbitalAuthentication = "orbital.authentication"
)

type LoginReq struct {
	RequestMessage string `json:"requestMessage,omitempty"`
}

type LoginResp struct {
	ConnectionID string `json:"connectionId"`
}

func (ws *WsConn) defaultOrbitalLogin(connID string, data []byte) {
	fmt.Println("Handle Login")
	var msg cryptographer.Message

	if err := json.Unmarshal(data, &msg); err != nil {
		ws.log.Error("error unmarshaling orbital login request:", err)
		return
	}

	var req LoginReq
	_ = json.Unmarshal(msg.Body, &req)

	fmt.Printf("MSG (income): %+v\n", req)

	//TODO: We need a way to sign the message. Maybe outsource this
	//TODO: TEST Data
	pk, sk, _ := cryptographer.GenerateKeysPair()

	res := &LoginResp{
		ConnectionID: connID,
	}

	resBytes, _ := json.Marshal(res)

	msgOut := NewTopicMessage(TopicOrbitalAuthentication, resBytes)
	msgOut.PublicKey = pk.Compress()
	msgOut.Sign(sk.Bytes())

	msgBytes, _ := json.Marshal(msgOut)

	pubKyStr := hex.EncodeToString(msg.PublicKey[:])
	
	ws.connectionManager.SetUserID(connID, pubKyStr)
	_ = ws.SendTo(connID, msgBytes)
}
