package auth

import (
	"context"
	"orbital/domain"
	"orbital/orbital"
	"orbital/pkg/logger"
)

const (
	Domain      = "auth"
	ActionLogin = "login"
)

type Dependencies struct {
	Log      *logger.Logger
	UserRepo domain.UserRepository
	Ws       *orbital.WsConn
}

type Auth struct {
	log *logger.Logger

	userRepo domain.UserRepository
	ws       *orbital.WsConn
}

func NewService(deps Dependencies) *Auth {
	return &Auth{
		log:      deps.Log,
		userRepo: deps.UserRepo,
		ws:       deps.Ws,
	}
}

func (service *Auth) Auth(ctx context.Context, req AuthReq) (*AuthResp, error) {

	userRepo, err := service.userRepo.FindByPublicKey(req.PublicKey)
	if err != nil {
		return &AuthResp{
			Code: orbital.NotFound,
			Error: &orbital.ErrorResponse{
				Type: "auth.notfound",
				Msg:  "unknown secret key",
			},
		}, nil
	}

	user := &User{
		ID:        userRepo.ID,
		Name:      userRepo.Name,
		PublicKey: userRepo.PubKey,
		Access:    userRepo.Access,
	}

	return &AuthResp{
		Code: orbital.OK,
		User: user,
	}, nil

}
