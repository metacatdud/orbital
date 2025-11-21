package auth

import (
	"context"
	"orbital/pkg/transport"
)

type AuthService interface {
	Auth(ctx context.Context, req AuthReq) (*AuthResp, error)
	Check(ctx context.Context, req CheckReq) (*CheckResp, error)
}

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
	Access    string `json:"access"`
}

type AuthReq struct {
	PublicKey string `json:"publicKey,omitempty"`
}

type AuthResp struct {
	User  *User                 `json:"user"`
	Code  transport.Code        `json:"code"`
	Error *transport.ErrorResponse `json:"error,omitempty"`
}

type CheckReq struct{}
type CheckResp struct {
	Code  transport.Code        `json:"code"`
	Error *transport.ErrorResponse `json:"error,omitempty"`
}
