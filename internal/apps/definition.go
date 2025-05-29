package apps

import (
	"context"
	"orbital/orbital"
)

type AppsService interface {
	List(ctx context.Context, req ListReq) (*ListResp, error)
}

type App struct {
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type ListReq struct{}

type ListResp struct {
	Code  orbital.Code           `json:"code"`
	Error *orbital.ErrorResponse `json:"error,omitempty"`
	Apps  []App                  `json:"apps"`
}
