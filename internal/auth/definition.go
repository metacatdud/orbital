package auth

import "context"

type AuthService interface {
	Auth(ctx context.Context, req AuthReq) (AuthResp, error)
}

type AuthReq struct {
	SecretKey string `json:"publicKey,omitempty"`
}

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
	Access    string `json:"access"`
}

type AuthResp struct {
	User  *User             `json:"user"`
	Error map[string]string `json:"error,omitempty"`
}
