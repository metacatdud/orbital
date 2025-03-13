package domain

import (
	"encoding/json"
	"orbital/web/wasm/pkg/transport"
)

const (
	UserStorageKey = "user"
)

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
	Access    string `json:"access"`
}

type LoginMessage struct {
	PublicKey string `json:"publicKey"`
}

type LoginResponse struct {
	Code  int                      `json:"code"`
	User  *User                    `json:"user"`
	Error *transport.ErrorResponse `json:"error,omitempty"`
}

func (msg *LoginResponse) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, msg)
}
