package orbital

import (
	"errors"
	"orbital/pkg/cryptographer"
	"orbital/web/wasm/components"
	"orbital/web/wasm/domain"
	"orbital/web/wasm/pkg/dom"
	"orbital/web/wasm/pkg/events"
	"orbital/web/wasm/pkg/state"
	"orbital/web/wasm/pkg/storage"
	"orbital/web/wasm/pkg/transport"
	"time"
)

type Orbital struct {
	di       *Dependency
	events   *events.Event
	state    *state.State
	rootComp *components.OrbitalComponent
	storage  storage.Storage
	ws       *transport.WsConn
}

func NewOrbital(di *Dependency) (*Orbital, error) {
	app := &Orbital{
		di:      di,
		events:  di.Events(),
		state:   di.State(),
		storage: di.Storage(),
		ws:      di.Ws(),
	}

	app.init()

	return app, nil
}

func (app *Orbital) Boot() {

	retries := 3
	interval := 1 * time.Second

	go func() {
		for i := 0; i < retries; i++ {
			if app.ws.IsOpen() {
				app.events.Emit("evt:app:ready")
				return
			}

			time.Sleep(interval)
		}

		app.state.Set("wsConnected", false)
		app.events.Emit("evt:app:ready")
	}()
}

func (app *Orbital) init() {
	app.events.Once("evt:app:ready", app.eventAppReady)
}

func (app *Orbital) eventAppReady() {
	dom.ConsoleLog("[orbital] Ready")

	rootEl := dom.QuerySelector("#rootEl")
	if rootEl.IsNull() {
		dom.ConsoleError("Element rootEl doesn't exist")
		return
	}

	app.rootComp = components.NewOrbitalComponent(app.di)
	if err := app.rootComp.Render(); err != nil {
		dom.ConsoleError("[orbital] Cannot render root component", err.Error())
		return
	}

	if err := app.rootComp.Mount(&rootEl); err != nil {
		dom.ConsoleError("[orbital] Cannot mount to rootEl", err.Error())
		return
	}

	authRepo := domain.NewRepository[domain.Auth](app.di.Storage(), domain.AuthStorageKey)
	auth, err := authRepo.Get()
	if err != nil {
		if !errors.Is(err, domain.ErrKeyNotFound) {
			dom.ConsoleError("[orbital] Cannot read storage", err.Error())
			return
		}
	}

	if auth == nil {
		app.state.Set("state:orbital:authenticated", false)
		return
	}

	userRepo := domain.NewRepository[domain.User](app.di.Storage(), domain.UserStorageKey)
	user, err := userRepo.Get()
	if err != nil {
		if !errors.Is(err, domain.ErrKeyNotFound) {
			dom.ConsoleError("[orbital] Cannot read storage", err.Error())
			return
		}
	}

	if user == nil {
		app.state.Set("state:orbital:authenticated", false)
		return
	}

	// Verify if stored key matches the user
	privateKey, err := cryptographer.NewPrivateKeyFromString(auth.SecretKey)
	if err != nil {
		dom.ConsoleError("[orbital] Cannot parse private key", err.Error())

		_ = authRepo.Remove()
		_ = userRepo.Remove()

		app.state.Set("state:orbital:authenticated", false)
		return
	}

	if privateKey.PublicKey().String() != user.PublicKey {
		dom.ConsoleError("[orbital] PrivateKey does not match public key")

		_ = authRepo.Remove()
		_ = userRepo.Remove()

		app.state.Set("state:orbital:authenticated", false)
		return
	}

	app.state.Set("state:orbital:authenticated", true)

}
