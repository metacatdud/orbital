package app

import (
	"fmt"
	"orbital/pkg/cryptographer"
	"orbital/web/wasm/api"
	"orbital/web/wasm/dom"
	"orbital/web/wasm/events"
	"orbital/web/wasm/storage"
	"syscall/js"
	"time"
)

type AppDI struct {
	Events  *events.Event
	Storage storage.Storage
	WsConn  *WsConn
}

type App struct {
	storage storage.Storage
	events  *events.Event
	wsConn  *WsConn
}

func (app *App) Boot() {
	app.events.Emit("app.ready")
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
		fmt.Println("orbital.authentication:", string(data))
	})
}

func (app *App) eventAppReady() {
	if app.hasSession() {
		var authData map[string]string
		if err := app.storage.Get("auth", &authData); err != nil {
			dom.PrintToConsole("Failed to get public key")
			return
		}

		//TODO: Validate to backend to
		fmt.Printf("Validate backend: %+v\n", authData)

		// Initialize the websocket connection
		app.events.Emit("dashboard.show")

		//TODO: MOVE THIS SOMEPLACE ELSE

		secretKy, err := cryptographer.NewPrivateKeyFromString("123")

		if err != nil {
			dom.PrintToConsole("Failed to convert Validate public key")
			return
		}

		async := api.NewAsync()
		async.Run(func() {
			time.Sleep(1 * time.Second)
			msg := NewTopicMessage("orbital.authentication", []byte(`{"requestMessage": "do.login"}`))
			msg.PublicKey = secretKy.PublicKey().Compress()
			if err = msg.Sign(secretKy.Bytes()); err != nil {
				dom.PrintToConsole("Failed to sign message")
				return
			}

			time.Sleep(1 * time.Second)
			fmt.Println("[HACK-ish] Wait for connections")
			app.wsConn.Send(*msg)
		})
		return
	}

	fmt.Println("Don't have session. Login")
	app.events.Emit("login.show")
}

func (app *App) eventAppNav(target string) {
	if target == "" {
		fmt.Println("Error: No target screen provided for navigation")
		return
	}

	fmt.Printf("Navigating to: %s\n", target)
	app.events.Emit(target + ".show")
}

func (app *App) hasSession() bool {
	var authData map[string]string
	err := app.storage.Get("auth", &authData)
	if err != nil {
		return false
	}

	if _, ok := authData["publicKey"]; !ok {
		return false
	}

	return true
}

func NewApp(di AppDI) *App {
	app := &App{
		events:  di.Events,
		storage: di.Storage,
		wsConn:  di.WsConn,
	}

	app.prepare()

	return app
}
