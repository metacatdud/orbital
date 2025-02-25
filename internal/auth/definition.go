package auth

import (
	"context"
	"orbital/orbital"
)

type AuthService interface {
	Auth(ctx context.Context, req AuthReq) (AuthResp, error)
	WsAuth(ctx context.Context, connID string, req WsAuthReq) error
}

type AuthReq struct {
	PublicKey string `json:"publicKey,omitempty"`
}

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
	Access    string `json:"access"`
}

type AuthResp struct {
	User  *User                  `json:"user"`
	Code  orbital.Code           `json:"code"`
	Error *orbital.ErrorResponse `json:"error,omitempty"`
}

type WsAuthReq struct {
	Authorize string `json:"authorize"`
}

type WsUser struct {
	ConnectionID string `json:"connectionId"`
	PublicKey    string `json:"publicKey"`
}

type WsAuthResp struct {
	User  *WsUser                `json:"user"`
	Code  orbital.Code           `json:"code"`
	Error *orbital.ErrorResponse `json:"error,omitempty"`
}
