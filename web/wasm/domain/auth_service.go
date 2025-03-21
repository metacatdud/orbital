package domain

import (
	"orbital/pkg/cryptographer"
	"orbital/pkg/proto"
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/events"
	"orbital/web/wasm/pkg/transport"
)

type AuthService struct {
	di     *orbital.Dependency
	events *events.Event
	ws     *transport.WsConn
}

func NewAuthService(di *orbital.Dependency) *AuthService {
	auth := &AuthService{
		di:     di,
		events: di.Events(),
		ws:     di.Ws(),
	}

	auth.init()

	return auth
}

func (auth *AuthService) init() {

	// Event listeners
	auth.events.On("evt:auth:login:request", auth.eventLogin)
	auth.events.On("evt:auth:logout:request", auth.eventLogout)
	auth.events.On("evt:auth:roleAccess:request", auth.eventRoleAccess)

	// Ws listeners
	auth.ws.On("orbital:authenticated", auth.wsAuthenticated)
}

func (auth *AuthService) eventLogin(secretKey string) {

	var loginErr *transport.ErrorResponse
	if secretKey == "" {
		loginErr = &transport.ErrorResponse{
			Type: "auth.empty",
			Msg:  "private key cannot be empty",
		}

		auth.events.Emit("evt:auth:login:fail", loginErr)
		return
	}

	sk, err := cryptographer.NewPrivateKeyFromString(secretKey)
	if err != nil {
		loginErr = &transport.ErrorResponse{
			Type: "auth.unauthorized",
			Msg:  "private key cannot be parsed",
		}
		auth.events.Emit("evt:auth:login:fail", loginErr)
		return
	}

	body := &LoginMessage{
		PublicKey: sk.PublicKey().String(),
	}

	meta := &transport.WsMetadata{
		Topic: "orbital:authenticate",
	}

	req, err := proto.Encode(sk, meta, body)
	if err != nil {
		loginErr = &transport.ErrorResponse{
			Type: "auth.unauthorized",
			Msg:  "private key not valid ed25519 key",
		}
		return
	}

	// If everything is OK we can save user private key right away
	authRepo := NewRepository[*Auth](auth.di.Storage(), AuthStorageKey)
	if err = authRepo.Save(&Auth{
		SecretKey: secretKey,
	}); err != nil {
		return
	}

	auth.ws.Send(*req)
}

func (auth *AuthService) eventLogout() {
	userRepo := NewRepository[*User](auth.di.Storage(), UserStorageKey)
	_ = userRepo.Remove()

	authRepo := NewRepository[*Auth](auth.di.Storage(), AuthStorageKey)
	_ = authRepo.Remove()

	auth.di.State().Set("state:orbital:authenticated", false)
}

func (auth *AuthService) eventRoleAccess(role string) {
	if role == "" {
		role = "guest"
	}

}

func (auth *AuthService) wsAuthenticated(data []byte) {
	loginRes := LoginResponse{}
	if err := loginRes.UnmarshalBinary(data); err != nil {
		auth.events.Emit("evt:auth:login:fail", &transport.ErrorResponse{
			Type: "auth.internal",
			Msg:  err.Error(),
		})
		return
	}

	userRepo := NewRepository[*User](auth.di.Storage(), UserStorageKey)
	if err := userRepo.Save(loginRes.User); err != nil {
		return
	}

	auth.events.Emit("evt:auth:login:success")
}
