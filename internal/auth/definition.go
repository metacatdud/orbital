package auth

import "context"

type AuthService interface {
	Auth(ctx context.Context, req AuthReq) (AuthResp, error)
}

type AuthReq struct {
	PublicKey string `json:"publicKey,omitempty"`
}

type AuthResp struct {
	Greet string `json:"greet,omitempty"`
}
