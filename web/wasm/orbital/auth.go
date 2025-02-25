package orbital

import (
	"encoding/json"
	"orbital/orbital"
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
}

func (auth *Auth) init() {
	auth.events.On("evt:auth:login:request", auth.eventLogin)
}

func (auth *Auth) eventLogin(secretKey string) {

	var loginErr *orbital.ErrorResponse
	if secretKey == "" {
		loginErr = &orbital.ErrorResponse{
			Type: "auth.empty",
			Msg:  "private key cannot be empty",
		}

		auth.events.Emit("auth::login.fail", loginErr)
		return
	}

	sk, err := cryptographer.NewPrivateKeyFromString(secretKey)
	if err != nil {
		loginErr = &orbital.ErrorResponse{
			Type: "auth.unauthorized",
			Msg:  "private key cannot be parsed",
		}
		auth.events.Emit("auth::login.fail", loginErr)
		return
	}

	loginReq := &domain.LoginMessage{
		PublicKey: sk.PublicKey().String(),
	}

	loginReqBin, err := loginReq.MarshalBinary()
	if err != nil {
		loginErr = &orbital.ErrorResponse{
			Type: "auth.unknown",
			Msg:  "cannot marshal login message",
		}
		auth.events.Emit("auth::login.fail", loginErr)
		return
	}

	req := &proto.Message{
		PublicKey: sk.PublicKey().Compress(),
		V:         1,
		Body:      loginReqBin,
		Timestamp: proto.TimestampNow(),
	}

	if err = req.Sign(sk.Bytes()); err != nil {
		loginErr = &orbital.ErrorResponse{
			Type: "auth.unauthorized",
			Msg:  "private key not valid ed25519 key",
		}

		auth.events.Emit("auth::login.fail", loginErr)
	}

	reqBin, err := json.Marshal(req)
	if err != nil {
		loginErr = &orbital.ErrorResponse{
			Type: "auth.unknown",
			Msg:  "cannot marshal request message",
		}
		auth.events.Emit("auth::login.fail", loginErr)
	}

	var async transport.Async
	async.Async(func() {
		client := transport.NewAPI("/rpc/AuthService/Auth")

		var res []byte
		res, err = client.Do(reqBin, nil)
		if err != nil {
			dom.ConsoleError("Error calling AuthService/Auth", err.Error())
			return
		}

		dom.ConsoleLog("AuthService/Auth success")

		loginRes := &domain.LoginResponse{}
		if err = loginRes.UnmarshalBinary(res); err != nil {
			dom.ConsoleLog("Error unmarshalling response", err.Error())
			return
		}

		if loginRes.Error != nil {
			loginErr = &orbital.ErrorResponse{
				Type: loginRes.Error.Type,
				Msg:  loginRes.Error.Msg,
			}

			auth.events.Emit("auth::login.fail", loginErr)
			return
		}

		userRepo := domain.NewRepository[*domain.User](auth.di.Storage(), domain.UserStorageKey)
		if err = userRepo.Save(loginRes.User); err != nil {
			dom.ConsoleLog("Error saving response", err.Error())
			return
		}

		authRepo := domain.NewRepository[*domain.Auth](auth.di.Storage(), domain.AuthStorageKey)
		err = authRepo.Save(&domain.Auth{
			SecretKey: secretKey,
		})
		if err != nil {
			dom.ConsoleLog("Error saving response", err.Error())
			return
		}

		auth.events.Emit("evt:auth:login:success")
		return
	})
}

func NewAuth(di *deps.Dependency) *Auth {
	auth := &Auth{
		di:     di,
		events: di.Events(),
	}

	auth.init()

	return auth
}
