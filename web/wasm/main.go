package main

import (
	"orbital/web/wasm/components"
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/dom"
	"orbital/web/wasm/pkg/transport"
	"orbital/web/wasm/service"
	"time"
)

func main() {

	boot()
	select {}
}

// boot load the entrypoint
func boot() {

	deps, err := orbital.NewDependency()
	if err != nil {
		dom.ConsoleLog("[orbital] cannot build dependencies")
		return
	}

	deps.Events.Once("orbital:ready", ready)

	retries := 3
	interval := 1 * time.Second

	var async transport.Async
	async.Async(func() {
		for i := 0; i < retries; i++ {
			if deps.Ws.IsOpen() {
				deps.Events.Emit("orbital:ready", deps)
				return
			}

			time.Sleep(interval)
		}

		dom.ConsoleError("[orbital] cannot open websocket connection")
	})
}

func ready(di *orbital.Dependency) {
	dom.ConsoleLog("[orbital] Ready")

	rootEl := dom.QuerySelector("#app-screen")
	if rootEl.IsNull() {
		dom.ConsoleError("[orbital] element rootEl doesn't exist")
		return
	}

	authSvc := service.NewAuthService(di)
	if err := di.RegisterService(service.AuthServiceKey, authSvc); err != nil {
		dom.ConsoleError("[orbital] cannot register service", service.AuthServiceKey)
	}

	appsSvc := service.NewAppsService(di)
	if err := di.RegisterService(service.AppsServiceKey, appsSvc); err != nil {
		dom.ConsoleError("[orbital] cannot register service", service.AppsServiceKey)
	}

	mainComp := components.NewMainComponent(di)
	_ = mainComp.Mount(&rootEl)

	checkAuthStatus(di)
}

func checkAuthStatus(di *orbital.Dependency) {

	authSvc := orbital.MustGetService[*service.AuthService](di, service.AuthServiceKey)

	var async transport.Async
	async.Async(func() {
		res, err := authSvc.CheckKey(service.CheckKeyReq{})
		if err != nil {
			di.State.Set("state:isAuthenticated", false)
			return
		}

		if res.Code == transport.OK {
			di.State.Set("state:isAuthenticated", true)
			return
		}

		di.State.Set("state:isAuthenticated", false)
	})

	async.Wait()
}
