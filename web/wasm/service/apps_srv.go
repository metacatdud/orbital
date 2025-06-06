package service

import (
	"encoding/json"
	"errors"
	"orbital/pkg/cryptographer"
	"orbital/pkg/proto"
	"orbital/web/wasm/domain"
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/transport"
)

const (
	AppsServiceKey = "appsServiceKey"
)

type App struct {
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Version     string `json:"version"`
	Description string `json:"description"`
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

	msg, err := proto.Encode(sk, &cryptographer.Metadata{
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
