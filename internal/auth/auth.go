package auth

import (
	"context"
	"orbital/domain"
	"orbital/orbital"
	"orbital/pkg/cryptographer"
)

type Dependencies struct {
	UserRepo domain.UserRepository
}

type Auth struct {
	userRepo domain.UserRepository
}

func (service *Auth) Auth(ctx context.Context, req AuthReq) (AuthResp, error) {

	sk, err := cryptographer.NewPrivateKeyFromString(req.SecretKey)
	if err != nil {
		return AuthResp{
			Code: orbital.InvalidRequest,
			Error: map[string]string{
				"auth.invalid": "invalid secret key",
			},
		}, nil
	}

	publicKey := sk.PublicKey()
	user, err := service.userRepo.FindByPublicKey(publicKey.String())
	if err != nil {
		return AuthResp{
			Code: orbital.NotFound,
			Error: map[string]string{
				"auth.notfound": "unknown secret key",
			},
		}, nil
	}

	return AuthResp{
		Code: orbital.OK,
		User: &User{
			ID:        user.ID,
			Name:      user.Name,
			PublicKey: publicKey.String(),
			Access:    user.Access,
		},
	}, nil
}

func NewService(deps Dependencies) *Auth {
	return &Auth{
		userRepo: deps.UserRepo,
	}
}
