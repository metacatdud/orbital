package orbital

import (
	"errors"
	"orbital/web/wasm/components"
	"orbital/web/wasm/domain"
	"orbital/web/wasm/pkg/deps"
	"orbital/web/wasm/pkg/dom"
	"orbital/web/wasm/pkg/events"
	"orbital/web/wasm/pkg/state"
	"orbital/web/wasm/pkg/storage"
	"orbital/web/wasm/pkg/transport"
	"time"
)

type App struct {
	di       *deps.Dependency
	events   *events.Event
	state    *state.State
	rootComp *components.OrbitalComponent
	storage  storage.Storage
	ws       *transport.WsConn
}

func (app *App) Boot() {

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

func (app *App) init() {
	app.events.Once("evt:app:ready", app.eventAppReady)
}

func (app *App) eventAppReady() {
	dom.ConsoleLog("[orbital] Ready")

	rootEl := dom.QuerySelector("#app")
	if rootEl.IsNull() {
		dom.ConsoleError("Element rootEl doesn't exist")
		return
	}

	app.rootComp = components.NewOrbitalComponent(app.di)

	if err := app.rootComp.Mount(&rootEl); err != nil {
		dom.ConsoleError("Cannot mount to rootEl", err.Error())
		return
	}

	authRepo := domain.NewRepository[*domain.Auth](app.di.Storage(), domain.AuthStorageKey)
	auth, err := authRepo.Get()
	if err != nil {
		if errors.Is(err, domain.ErrKeyNotFound) {
			app.state.Set("state:orbital:authenticated", false)
		}
	}

	if auth != nil {
		app.state.Set("state:orbital:authenticated", true)
	}

	return

}

func NewApp(di *deps.Dependency) (*App, error) {

	app := &App{
		di:      di,
		events:  di.Events(),
		state:   di.State(),
		storage: di.Storage(),
		ws:      di.Ws(),
	}

	app.init()

	return app, nil
}
