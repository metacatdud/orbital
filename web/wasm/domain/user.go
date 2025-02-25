package domain

import (
	"encoding/json"
	"orbital/orbital"
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

func (msg *LoginMessage) MarshalBinary() ([]byte, error) {
	return json.Marshal(msg)
}

type LoginMetadata struct {
}

func (msg *LoginMetadata) MarshalBinary() ([]byte, error) {
	return json.Marshal(msg)
}

type LoginResponse struct {
	User  *User                  `json:"user"`
	Error *orbital.ErrorResponse `json:"error,omitempty"`
}

func (msg *LoginResponse) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, msg)
}
