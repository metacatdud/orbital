package orbital

import (
	"orbital/pkg/cryptographer"
	"orbital/pkg/proto"
	"orbital/web/wasm/domain"
	"orbital/web/wasm/pkg/deps"
	"orbital/web/wasm/pkg/dom"
	"orbital/web/wasm/pkg/events"
	"orbital/web/wasm/pkg/transport"
)

type Auth struct {
	di     *deps.Dependency
	events *events.Event
	ws     *transport.WsConn
}

func NewAuth(di *deps.Dependency) *Auth {
	auth := &Auth{
		di:     di,
		events: di.Events(),
		ws:     di.Ws(),
	}

	auth.init()

	return auth
}

func (auth *Auth) init() {

	// Event listeners
	auth.events.On("evt:auth:login:request", auth.eventLogin)
	auth.events.On("evt:auth:roleAccess:request", auth.eventRoleAccess)

	// Ws listeners
	auth.ws.On("orbital:authenticated", auth.wsAuthenticated)
}

func (auth *Auth) eventLogin(secretKey string) {

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

	body := &domain.LoginMessage{
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
	authRepo := domain.NewRepository[*domain.Auth](auth.di.Storage(), domain.AuthStorageKey)
	if err = authRepo.Save(&domain.Auth{
		SecretKey: secretKey,
	}); err != nil {
		dom.ConsoleLog("Error saving response", err.Error())
		return
	}

	auth.ws.Send(*req)
}

func (auth *Auth) eventRoleAccess(role string) {
	if role == "" {
		role = "guest"
	}

}

func (auth *Auth) wsAuthenticated(data []byte) {
	loginRes := domain.LoginResponse{}
	if err := loginRes.UnmarshalBinary(data); err != nil {
		auth.events.Emit("evt:auth:login:fail", &transport.ErrorResponse{
			Type: "auth.internal",
			Msg:  err.Error(),
		})
		return
	}

	userRepo := domain.NewRepository[*domain.User](auth.di.Storage(), domain.UserStorageKey)
	if err := userRepo.Save(loginRes.User); err != nil {
		dom.ConsoleLog("Error saving response", err.Error())
		return
	}

	auth.events.Emit("evt:auth:login:success")
}
