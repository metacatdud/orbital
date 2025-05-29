package main

import (
	"orbital/web/wasm/components"
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/dom"
	"orbital/web/wasm/pkg/events"
	"orbital/web/wasm/pkg/transport"
	"orbital/web/wasm/service"
	"syscall/js"
	"time"
)

func main() {

	js.Global().Set("bootstrapApp", js.FuncOf(bootstrapApp))

	select {}
}

// bootstrapApp load the entrypoint
func bootstrapApp(_ js.Value, _ []js.Value) interface{} {

	deps, err := orbital.NewDependency()
	if err != nil {
		dom.ConsoleLog("Cannot build dependencies")
		return nil
	}

	deps.Events.Once("orbital:ready", func() {
		ready(deps)
	})

	wsStatusCheck(deps.Ws, deps.Events)

	return nil
}

func ready(di *orbital.Dependency) {
	dom.ConsoleLog("[orbital] Ready")

	rootEl := dom.QuerySelector("#rootEl")
	if rootEl.IsNull() {
		dom.ConsoleError("Element rootEl doesn't exist")
		return
	}

	authSvc := service.NewAuthService(di)
	if err := di.RegisterService(service.AuthServiceKey, authSvc); err != nil {
		dom.ConsoleError("Cannot register auth service")
	}

	mainComp := components.NewMainComponent(di)
	_ = mainComp.Mount(&rootEl)

}

func wsStatusCheck(ws *transport.WsConn, evt *events.Event) {
	retries := 3
	interval := 1 * time.Second

	var async transport.Async
	async.Async(func() {
		for i := 0; i < retries; i++ {
			if ws.IsOpen() {
				evt.Emit("orbital:ready")
				return
			}

			time.Sleep(interval)
		}
	})
}
