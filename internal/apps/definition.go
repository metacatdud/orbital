package apps

import (
	"context"
	"orbital/pkg/transport"
)

type AppsService interface {
	List(ctx context.Context, req ListReq) (*ListResp, error)
}

// App struct holder for apps
// TODO: This needs extra properties
type App struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Icon        string   `json:"icon"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Namespace   string   `json:"namespace"`
	OwnerKey    string   `json:"ownerKey"`
	OwnerURL    string   `json:"ownerUrl"`
	Labels      []string `json:"labels"`
	Apps        []App    `json:"apps"` // If an app it's a suite of apps (just a group basically)
}

type ListReq struct{}

type ListResp struct {
	Code  transport.Code          `json:"code"`
	Error *transport.ErrorResponse `json:"error,omitempty"`
	Apps  []App                   `json:"apps"`
}
