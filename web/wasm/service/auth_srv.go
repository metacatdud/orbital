package service

import (
	"encoding/json"
	"orbital/pkg/cryptographer"
	"orbital/pkg/proto"
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

type LoginReq struct {
	SecretKey string `json:"secretKey"`
}

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Access string `json:"access"`
}

type LoginRes struct {
	Code  int                      `json:"code"`
	User  *User                    `json:"user"`
	Error *transport.ErrorResponse `json:"error,omitempty"`
}

func (srv *AuthService) Login(req LoginReq) (*LoginRes, error) {

	api := transport.NewAPI("rpc/AuthService/Auth")
	api.WithMiddleware(transport.VerifyAndUnwrap)

	sk, err := cryptographer.NewPrivateKeyFromString(req.SecretKey)
	if err != nil {
		return nil, err
	}

	msg, err := proto.Encode(sk, &cryptographer.Metadata{
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
