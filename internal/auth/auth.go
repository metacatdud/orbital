package auth

import (
	"context"
	"orbital/domain"
	"orbital/pkg/logger"
	"orbital/pkg/transport"

	"atomika.io/atomika/atomika"
)

const (
	Domain      = "auth"
	ActionLogin = "login"
	ActionCheck = "check"
)

type Dependencies struct {
	Log      *logger.Logger
	UserRepo *domain.UserRepository
	Ws       atomika.WSDispatcher
}

type Auth struct {
	log *logger.Logger

	userRepo *domain.UserRepository
	ws       atomika.WSDispatcher
}

func NewService(deps Dependencies) *Auth {
	return &Auth{
		log:      deps.Log,
		userRepo: deps.UserRepo,
		ws:       deps.Ws,
	}
}

func (service *Auth) Auth(ctx context.Context, req AuthReq) (*AuthResp, error) {

	userRepo, err := service.userRepo.GetByPublicKey(req.PublicKey)
	if err != nil {
		return &AuthResp{
			Code: transport.NotFound,
			Error: &transport.ErrorResponse{
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
		Code: transport.OK,
		User: user,
	}, nil

}

func (service *Auth) Check(ctx context.Context, req CheckReq) (*CheckResp, error) {
	// TODO :Check with database
	return &CheckResp{
		Code: transport.OK,
	}, nil

}
