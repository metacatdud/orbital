package app

import (
	"fmt"
	"orbital/pkg/cryptographer"
	"orbital/pkg/proto"
	"orbital/web/wasm/api"
	"orbital/web/wasm/dom"
	"orbital/web/wasm/domain"
	"orbital/web/wasm/events"
	"syscall/js"
	"time"
)

type AppDI struct {
	Events   *events.Event
	WsConn   *api.WsConn
	AuthRepo domain.AuthRepository
	UserRepo domain.UserRepository
}

type App struct {
	events   *events.Event
	wsConn   *api.WsConn
	authRepo domain.AuthRepository
	userRepo domain.UserRepository
}

func (app *App) Boot() {

	retries := 3
	interval := 1 * time.Second

	go func() {
		for i := 0; i < retries; i++ {
			if app.wsConn.IsOpen() {
				app.events.Emit("app.ready")
				return
			}

			time.Sleep(interval)
		}

		dom.PrintToConsole("Cannot acquire readiness checker")
	}()
}

func (app *App) Render(htmlEl js.Value) {
	dom.Clear("#app")

	if htmlEl.Truthy() {
		dom.AppendChild("#app", htmlEl)
	}

	dom.Hide("#loading")
	dom.Show("#app")
}

func (app *App) prepare() {
	app.events.Once("app.ready", app.eventAppReady)
	app.events.On("navigate", app.eventAppNav)
	app.events.On("app.render", func(tpl js.Value) {
		app.Render(tpl)
	})

	app.wsConn.On("orbital.authentication", func(data []byte) {
		dom.PrintToConsole("orbital.authentication:", string(data))
	})
}

func (app *App) eventAppReady() {
	if !app.userRepo.HasSession() {
		app.events.Emit("login.show")
	}

	auth, err := app.authRepo.Get()
	if err != nil {
		dom.PrintToConsole("Failed to get public key")
		return
	}

	app.events.Emit("dashboard.show")

	//TODO: MOVE THIS SOMEPLACE ELSE

	secretKey, err := cryptographer.NewPrivateKeyFromString(auth.SecretKey)
	if err != nil {
		dom.PrintToConsole("Failed to convert Validate public key")
		return
	}

	async := api.NewAsync()
	async.Run(func() {

		meta := &api.WsMetadata{
			Topic: "orbital.authentication",
		}

		body := map[string]string{
			"authorize": secretKey.PublicKey().String(),
		}

		var m *proto.Message
		m, err = proto.Encode(*secretKey, meta, body)
		if err != nil {
			dom.PrintToConsole("Failed to encode message")
		}
		
		app.wsConn.Send(*m)
	})
	return

}

func (app *App) eventAppNav(target string) {
	if target == "" {
		fmt.Println("Error: No target screen provided for navigation")
		return
	}

	fmt.Printf("Navigating to: %s\n", target)
	app.events.Emit(target + ".show")
}

func NewApp(di AppDI) *App {
	app := &App{
		events:   di.Events,
		wsConn:   di.WsConn,
		authRepo: di.AuthRepo,
		userRepo: di.UserRepo,
	}

	app.prepare()

	return app
}
