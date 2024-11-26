package app

import (
	"encoding/json"
	"fmt"
	"orbital/dashboard/wasm/dom"
	"orbital/dashboard/wasm/events"
	"orbital/dashboard/wasm/storage"
	"syscall/js"
)

type AppDI struct {
	Events  *events.Event
	Storage storage.Storage
}

type App struct {
	storage storage.Storage
	events  *events.Event
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
	app.events.On("app.ready", app.eventAppReady)
	app.events.On("navigate", app.eventAppNav)

	app.events.On("app.render", func(tpl js.Value) {
		fmt.Println("render node")
		app.Render(tpl)
	})
}

func (app *App) eventAppReady() {
	if app.hasSession() {
		var pubKey string
		if err := app.storage.Get("publicKey", &pubKey); err != nil {
			dom.PrintToConsole("Failed to get public key")
			return
		}

		data, err := json.Marshal(map[string]string{"pubKey": pubKey})
		if err != nil {
			dom.PrintToConsole("Failed to marshal public key")
			return
		}

		//TODO: Validate to backend to
		fmt.Printf("Validate backend: %+v\n", data)

		app.events.Emit("navigate", "dashboard.show")
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
	var publicKey string
	err := app.storage.Get("publicKey", &publicKey)
	if err != nil {
		return false
	}

	if publicKey == "" {
		return false
	}

	return true
}

func NewApp(di AppDI) *App {
	app := &App{
		events:  di.Events,
		storage: di.Storage,
	}

	app.prepare()

	return app
}
