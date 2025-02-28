package main

import (
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/deps"
	"orbital/web/wasm/pkg/dom"
	"orbital/web/wasm/pkg/events"
	"orbital/web/wasm/pkg/state"
	"orbital/web/wasm/pkg/storage"
	"orbital/web/wasm/pkg/transport"
	"orbital/web/wasm/templates"
	"syscall/js"
)

func main() {

	js.Global().Set("bootstrapApp", js.FuncOf(bootstrapApp))

	select {}
}

// bootstrapApp load the entrypoint
func bootstrapApp(_ js.Value, _ []js.Value) interface{} {

	tplReg, err := templates.NewRegistry()
	if err != nil {
		dom.ConsoleError("Cannot create template registry", err.Error())
		return nil
	}

	di, err := deps.NewDependency(deps.Packages{
		Events:      events.New(),
		State:       state.New(),
		Storage:     storage.NewLocalStorage(),
		TplRegistry: tplReg,
		Ws:          transport.NewWsConn(true),
	})

	if err != nil {
		dom.ConsoleError("Cannot create dependencies registry", err.Error())
		return nil
	}

	// Services
	_ = orbital.NewAuth(di)
	_ = orbital.NewMachine(di)

	// Boot Orbital
	app, err := orbital.NewApp(di)
	if err != nil {
		dom.ConsoleError("Cannot create orbital app", err.Error())
		return nil
	}

	app.Boot()

	return nil
}
