package apps

import (
	"context"
	"orbital/orbital"
)

type AppsService interface {
	List(ctx context.Context, req ListReq) (*ListResp, error)
}

// App struct holder for apps
// TODO: This needs extra properties
type App struct {
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Apps        []App  `json:"apps"` // If an app it's a suite of apps (just a group basically)
}

type ListReq struct{}

type ListResp struct {
	Code  orbital.Code           `json:"code"`
	Error *orbital.ErrorResponse `json:"error,omitempty"`
	Apps  []App                  `json:"apps"`
}
