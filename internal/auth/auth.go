package auth

import (
	"context"
	"fmt"
	"orbital/domain"
	"orbital/orbital"
	"orbital/pkg/cryptographer"
	"orbital/pkg/proto"
)

type Dependencies struct {
	UserRepo domain.UserRepository
	Ws       *orbital.WsConn
}

type Auth struct {
	userRepo domain.UserRepository
	ws       *orbital.WsConn
}

func (service *Auth) Auth(ctx context.Context, req AuthReq) (AuthResp, error) {

	user, err := service.userRepo.FindByPublicKey(req.PublicKey)
	if err != nil {
		return AuthResp{
			Code: orbital.NotFound,
			Error: &orbital.ErrorResponse{
				Type: "auth.notfound",
				Msg:  "unknown secret key",
			},
		}, nil
	}

	return AuthResp{
		Code: orbital.OK,
		User: &User{
			ID:        user.ID,
			Name:      user.Name,
			PublicKey: req.PublicKey,
			Access:    user.Access,
		},
	}, nil
}

func (service *Auth) WsAuth(ctx context.Context, connID string, req WsAuthReq) error {
	fmt.Printf("Authorize WS sock connection for: %s\n", req.Authorize)

	// DUMMY SERVER KEYS
	_, sk, _ := cryptographer.GenerateKeysPair()

	body := &WsAuthResp{}
	user, err := service.userRepo.FindByPublicKey(req.Authorize)
	if err != nil {

		meta := &orbital.WsMetadata{
			Topic: "orbital.authenticationFail",
		}

		body.Code = orbital.Unauthenticated
		body.Error = &orbital.ErrorResponse{
			Type: "auth.unauthenticated",
			Msg:  "unknown secret key",
		}

		msg, _ := proto.Encode(*sk, meta, body)
		if err = service.ws.SendTo(connID, *msg); err != nil {
			return nil
		}
	}

	meta := &orbital.WsMetadata{
		Topic: "orbital.authenticationSuccess",
	}

	resUser := &WsUser{
		ConnectionID: connID,
		PublicKey:    user.PubKey,
	}

	body.User = resUser
	body.Code = orbital.OK

	msg, _ := proto.Encode(*sk, meta, body)
	if err = service.ws.SendTo(connID, *msg); err != nil {
		return nil
	}

	return nil
}

func NewService(deps Dependencies) *Auth {
	return &Auth{
		userRepo: deps.UserRepo,
		ws:       deps.Ws,
	}
}
