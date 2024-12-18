package app

import (
	"fmt"
	"orbital/web/wasm/dom"
	"orbital/web/wasm/events"
	"orbital/web/wasm/storage"
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
		app.Render(tpl)
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

		app.events.Emit("dashboard.show")
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
	}

	app.prepare()

	return app
}
