package service

import "orbital/web/wasm/orbital"

const (
	AppsServiceKey = "appsServiceKey"
)

type AppsService struct {
	di *orbital.Dependency
}

func NewSystemService(di *orbital.Dependency) *AppsService {
	return &AppsService{
		di: di,
	}
}

func (s *AppsService) ID() string {
	return AppsServiceKey
}

func (s *AppsService) List() {
	// TODO: call rpc/AppsService/List here
}
