package service

import "orbital/web/w2/orbital"

type AuthService struct {
	di *orbital.Dependency
}

func NewAuthService(di *orbital.Dependency) *AuthService {
	return &AuthService{
		di: di,
	}
}
