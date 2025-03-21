package auth

import (
	"context"
	"orbital/domain"
	"orbital/orbital"
	"orbital/pkg/cryptographer"
	"orbital/pkg/logger"
	"orbital/pkg/proto"
)

type Dependencies struct {
	Log      *logger.Logger
	NodePk   *cryptographer.PrivateKey
	UserRepo domain.UserRepository
	Ws       *orbital.WsConn
}

type Auth struct {
	log      *logger.Logger
	nodePk   *cryptographer.PrivateKey
	userRepo domain.UserRepository
	ws       *orbital.WsConn
}

func NewService(deps Dependencies) *Auth {
	return &Auth{
		log:      deps.Log,
		nodePk:   deps.NodePk,
		userRepo: deps.UserRepo,
		ws:       deps.Ws,
	}
}

func (service *Auth) Auth(ctx context.Context, connID string, req AuthReq) error {

	userRepo, err := service.userRepo.FindByPublicKey(req.PublicKey)
	meta := &orbital.WsMetadata{
		Topic: "orbital:authenticated",
	}

	if err != nil {
		body := &AuthResp{
			Code: orbital.NotFound,
			Error: &orbital.ErrorResponse{
				Type: "auth.notfound",
				Msg:  "unknown secret key",
			},
		}

		msg, _ := proto.Encode(service.nodePk, meta, body)
		if err = service.ws.SendTo(connID, *msg); err != nil {
			service.log.Error("send message failed", "err", err.Error(), "connID", connID)
			return nil
		}

		return nil
	}

	user := &User{
		ID:        userRepo.ID,
		Name:      userRepo.Name,
		PublicKey: userRepo.PubKey,
		Access:    userRepo.Access,
	}

	body := &AuthResp{
		Code: orbital.OK,
		User: user,
	}

	msg, _ := proto.Encode(service.nodePk, meta, body)
	if err = service.ws.SendTo(connID, *msg); err != nil {
		service.log.Error("send message failed", "err", err.Error(), "connID", connID)
		return nil
	}

	return nil

}
