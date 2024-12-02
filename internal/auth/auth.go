package auth

import (
	"context"
	"fmt"
)

type Hello struct {
}

func (service *Hello) Auth(ctx context.Context, req AuthReq) (AuthResp, error) {
	return AuthResp{
		Greet: fmt.Sprintf("Hello :%s", req.PublicKey),
	}, nil
}

func NewService() *Hello {
	return &Hello{}
}
