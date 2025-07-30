package apps

import (
	"context"
	"fmt"
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

	dbApps, err = service.appRepo.FindOnlyStandalone()
	if err != nil {
		return nil, err
	}

	var apps []App
	for _, dbApp := range dbApps {
		var appsTree App

		appsTree, err = service.buildTree(dbApp)
		if err != nil {
			return nil, err
		}

		apps = append(apps, appsTree)
	}

	return &ListResp{
		Code: orbital.OK,
		Apps: apps,
	}, nil
}

// buildTree - recursive apps retrieval
func (service *Apps) buildTree(app domain.App) (App, error) {
	children, err := service.appRepo.FindByParentID(app.ID)
	if err != nil {
		return App{}, fmt.Errorf("error fetching children for %s: %w", app.ID, err)
	}

	var childrenList []App
	for _, child := range children {
		var childNode App
		childNode, err = service.buildTree(child)
		if err != nil {
			return App{}, err
		}

		childrenList = append(childrenList, childNode)
	}

	return App{
		ID:          app.ID,
		Name:        app.Name,
		Icon:        app.Icon,
		Version:     app.Version,
		Description: app.Description,
		Namespace:   app.Namespace,
		OwnerKey:    app.OwnerKey,
		OwnerURL:    app.OwnerURL,
		Labels:      app.Labels,
		Apps:        childrenList,
	}, nil
}
