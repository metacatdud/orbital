package service

import (
	"encoding/json"
	"errors"
	"orbital/pkg/cryptographer"
	"orbital/web/wasm/domain"
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/transport"
)

const (
	AppsServiceKey = "appsServiceKey"
)

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
	IsExternal  bool     `json:"isExternal"`
	Apps        []App    `json:"apps"`
}

type AppsService struct {
	di *orbital.Dependency
}

func NewAppsService(di *orbital.Dependency) *AppsService {
	return &AppsService{
		di: di,
	}
}

func (srv *AppsService) ID() string {
	return AppsServiceKey
}

func (srv *AppsService) List(req ListReq) (*ListRes, error) {
	authRepo := domain.NewAuthRepository(srv.di.Storage)
	auth, err := authRepo.Get()
	if err != nil {
		if errors.Is(err, domain.ErrKeyNotFound) {
			return nil, domain.ErrKeyNotFound
		}

		return nil, err
	}

	sk, err := cryptographer.NewPrivateKeyFromString(auth.SecretKey)
	if err != nil {
		return nil, err
	}

	api := transport.NewAPI("rpc/AppsService/List")
	api.WithMiddleware(transport.VerifyAndUnwrap)

	msg, err := cryptographer.Encode(sk, &cryptographer.Metadata{
		Domain: "apps",
		Action: "list",
	}, req)

	raw, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	var (
		res    *ListRes
		rawRes []byte
	)

	rawRes, err = api.Do(raw, nil)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(rawRes, &res); err != nil {
		return nil, err
	}

	if res.Error != nil {
		return res, nil
	}

	return res, nil
}

type (
	// ListReq TODO: enrich with filter if needed in the future
	ListReq struct{}
	ListRes struct {
		Code  int                      `json:"code"`
		Error *transport.ErrorResponse `json:"error,omitempty"`
		Apps  []App                    `json:"apps"`
	}
)
