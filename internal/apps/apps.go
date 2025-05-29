package apps

import (
	"context"
	"orbital/orbital"
	"orbital/pkg/logger"
)

const (
	Domain     = "apps"
	ActionList = "list"
)

type Dependencies struct {
	Log *logger.Logger
}

type Apps struct {
	log *logger.Logger
}

func NewService(deps Dependencies) *Apps {
	return &Apps{
		log: deps.Log,
	}
}

func (a *Apps) List(ctx context.Context, req ListReq) (*ListResp, error) {

	// These should come from DB
	dummyApps := []App{
		{
			Name:        "Dummy App 1",
			Icon:        "fa-brands fa-medapps",
			Version:     "1.0.0",
			Description: "This is my first app",
		},
		{
			Name:        "Notes",
			Icon:        "fa-note-sticky",
			Version:     "1.1.0",
			Description: "Notes to read never",
		},
		{
			Name:        "System",
			Icon:        "fa-gear",
			Version:     "2.0.10",
			Description: "Orbital settings",
		},
	}

	return &ListResp{
		Code: orbital.Unimplemented,
		Apps: dummyApps,
	}, nil
}
