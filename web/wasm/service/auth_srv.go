package service

import (
	"encoding/json"
	"errors"
	"orbital/pkg/cryptographer"
	"orbital/web/wasm/domain"
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/events"
	"orbital/web/wasm/pkg/transport"
)

const (
	AuthServiceKey = "authService"
)

type AuthService struct {
	di *orbital.Dependency
}

func NewAuthService(di *orbital.Dependency) *AuthService {
	return &AuthService{
		di: di,
	}
}

func (srv *AuthService) ID() string {
	return AuthServiceKey
}

// HookEvents register events for this service
func (srv *AuthService) HookEvents(ev *events.Event) {
	ev.On("auth:login", srv.Login)
}

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Access string `json:"access"`
}

func (srv *AuthService) Login(req LoginReq) (*LoginRes, error) {

	api := transport.NewAPI("rpc/AuthService/Auth")
	api.WithMiddleware(transport.VerifyAndUnwrap)

	sk, err := cryptographer.NewPrivateKeyFromString(req.SecretKey)
	if err != nil {
		return nil, err
	}

	msg, err := cryptographer.Encode(sk, &cryptographer.Metadata{
		Domain: "auth",
		Action: "login",
	}, map[string]any{
		"publicKey": sk.PublicKey().String(),
	})

	raw, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	var (
		res    *LoginRes
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

	// Store secret key localstorage (A browser extension might be a better approach for future)
	authRepo := domain.NewAuthRepository(srv.di.Storage)
	if err = authRepo.Save(domain.Auth{
		SecretKey: req.SecretKey,
	}); err != nil {
		return nil, err
	}

	userRepo := domain.NewUserRepository(srv.di.Storage)
	err = userRepo.Save(domain.User{
		ID:     res.User.ID,
		Name:   res.User.Name,
		Access: res.User.Access,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

type (
	LoginReq struct {
		SecretKey string `json:"secretKey"`
	}

	LoginRes struct {
		Code  transport.Code           `json:"code"`
		User  *User                    `json:"user"`
		Error *transport.ErrorResponse `json:"error,omitempty"`
	}
)

func (srv *AuthService) CheckKey(req CheckKeyReq) (*CheckKeyRes, error) {
	authRepo := domain.NewAuthRepository(srv.di.Storage)
	auth, err := authRepo.Get()
	if err != nil {
		if errors.Is(err, domain.ErrKeyNotFound) {
			return &CheckKeyRes{Code: transport.Unauthenticated}, nil
		}

		return nil, err
	}

	sk, err := cryptographer.NewPrivateKeyFromString(auth.SecretKey)
	if err != nil {
		return &CheckKeyRes{Code: transport.Unauthenticated}, nil
	}

	api := transport.NewAPI("rpc/AuthService/Check")
	api.WithMiddleware(transport.VerifyAndUnwrap)

	msg, err := cryptographer.Encode(sk, &cryptographer.Metadata{
		Domain: "auth",
		Action: "check",
	}, nil)

	raw, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	var (
		res    *CheckKeyRes
		rawRes []byte
	)
	rawRes, err = api.Do(raw, nil)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(rawRes, &res); err != nil {
		return nil, err
	}

	return &CheckKeyRes{Code: res.Code}, nil
}

type (
	CheckKeyReq struct{}
	CheckKeyRes struct {
		Code  transport.Code           `json:"code"`
		Error *transport.ErrorResponse `json:"error,omitempty"`
	}
)
