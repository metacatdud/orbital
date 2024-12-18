package auth

import (
	"context"
	"orbital/domain"
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
			Error: map[string]string{
				"auth.invalid": "invalid secret key",
			},
		}, nil
	}

	publickKey := sk.PublicKey()
	user, err := service.userRepo.FindByPublicKey(publickKey.String())
	if err != nil {
		return AuthResp{
			Error: map[string]string{
				"auth.notfound": "unknown secret key",
			},
		}, nil
	}

	return AuthResp{
		User: &User{
			ID:        user.ID,
			Name:      user.Name,
			PublicKey: publickKey.String(),
			Access:    user.Access,
		},
	}, nil
}

func NewService(deps Dependencies) *Auth {
	return &Auth{
		userRepo: deps.UserRepo,
	}
}
