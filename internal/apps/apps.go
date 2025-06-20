package apps

import (
	"context"
	"orbital/domain"
	"orbital/orbital"
	"orbital/pkg/logger"
)

const (
	Domain     = "apps"
	ActionList = "list"
)

type Dependencies struct {
	Log     *logger.Logger
	AppRepo *domain.AppRepository
}

type Apps struct {
	log     *logger.Logger
	appRepo *domain.AppRepository
}

func NewService(deps Dependencies) *Apps {
	return &Apps{
		log:     deps.Log,
		appRepo: deps.AppRepo,
	}
}

func (service *Apps) List(_ context.Context, _ ListReq) (*ListResp, error) {

	var (
		dbApps domain.Apps
		err    error
	)

	dbApps, err = service.appRepo.Find()
	if err != nil {
		return nil, err
	}

	var apps []App
	for _, dbApp := range dbApps {
		apps = append(apps, App{
			Name:        dbApp.Name,
			Icon:        dbApp.Icon,
			Version:     dbApp.Version,
			Description: dbApp.Description,
			Apps:        nil,
		})
	}

	return &ListResp{
		Code: orbital.Unimplemented,
		Apps: apps,
	}, nil
}
