package main

import (
	"orbital/web/w2/components"
	"orbital/web/w2/orbital"
	"orbital/web/w2/pkg/dom"
	"orbital/web/w2/pkg/events"
	"orbital/web/w2/pkg/transport"
	"orbital/web/w2/service"
	"syscall/js"
	"time"
)

func main() {

	js.Global().Set("bootstrapApp", js.FuncOf(bootstrapApp))

	select {}
}

// bootstrapApp load the entrypoint
func bootstrapApp(_ js.Value, _ []js.Value) interface{} {

	deps, err := orbital.NewDependencyWithDefaults()
	if err != nil {
		dom.ConsoleLog("Cannot build dependencies")
		return nil
	}

	deps.Events().Once("orbital:ready", ready)

	orbital.NewRegistry(deps)

	service.Init()
	components.Init()

	wsStatusCheck(deps.Ws(), deps.Events())

	return nil
}

func ready() {
	dom.ConsoleLog("[orbital] Ready")

	rootEl := dom.QuerySelector("#rootEl")
	if rootEl.IsNull() {
		dom.ConsoleError("Element rootEl doesn't exist")
		return
	}

	orbitalComp, err := orbital.Lookup[*components.OrbitalComp](components.OrbitalCompKey)
	if err != nil {
		dom.ConsoleError("[orbital] Cannot create OrbitalComp", err.Error())
		return
	}

	if err = orbitalComp.Render(); err != nil {
		dom.ConsoleError("[orbital] Cannot render OrbitalComp", err.Error())
		return
	}

	if err = orbitalComp.Mount(&rootEl); err != nil {
		dom.ConsoleError("[orbital] Cannot mount OrbitalComp", err.Error())
		return
	}
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

//func setupFactory(di *orbital.Dependency) *orbital.Factory {
//	factory := orbital.NewFactory(di)
//
//	factory.Register(components.OrbitalCompKey, func(di *orbital.Dependency, params ...interface{}) orbital.Component {
//		return components.NewOrbitalComp(di, factory)
//	})
//
//	factory.Register(components.TaskbarCompKey, func(di *orbital.Dependency, params ...interface{}) orbital.Component {
//		return components.NewTaskbarComp(di, factory)
//	})
//
//	factory.Register(components.TaskbarStartCompKey, func(di *orbital.Dependency, params ...interface{}) orbital.Component {
//		return components.NewTaskbarStartComp(di, factory)
//	})
//
//	return factory
//}
