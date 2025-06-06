package auth

import (
	"context"
	"orbital/orbital"
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
	User  *User                  `json:"user"`
	Code  orbital.Code           `json:"code"`
	Error *orbital.ErrorResponse `json:"error,omitempty"`
}

type CheckReq struct{}
type CheckResp struct {
	Code  orbital.Code           `json:"code"`
	Error *orbital.ErrorResponse `json:"error,omitempty"`
}
