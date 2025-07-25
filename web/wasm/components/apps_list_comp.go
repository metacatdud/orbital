package components

import (
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/dom"
	"orbital/web/wasm/pkg/state"
	"orbital/web/wasm/pkg/transport"
	"orbital/web/wasm/service"
)

const AppsListComponentRegKey RegKey = "appsListComponent"

type AppsListComponent struct {
	*BaseComponent
	DI      *orbital.Dependency
	state   *state.State
	appsSvc *service.AppsService
	apps    []service.App
}

func NewAppsListComponent(di *orbital.Dependency) *AppsListComponent {
	base := NewBaseComponent(di, AppsListComponentRegKey, "dashboard/app/applications")

	// Services
	appsSvc := orbital.MustGetService[*service.AppsService](di, service.AppsServiceKey)

	comp := &AppsListComponent{
		BaseComponent: base,
		DI:            di,
		state:         di.State,
		appsSvc:       appsSvc,
	}

	comp.OnMount(func() {
		comp.loadApps()
	})

	comp.state.Watch("state:apps:changed", func(_, _ any) {
		comp.updateApps()
	})

	return comp
}

func (comp *AppsListComponent) ID() RegKey {
	return AppsListComponentRegKey
}

func (comp *AppsListComponent) loadApps() {
	var async transport.Async
	async.Async(func() {
		res, err := comp.appsSvc.List(service.ListReq{})
		if err != nil {
			dom.ConsoleError(err)
			return
		}

		if res.Error != nil {
			dom.ConsoleError(res.Error)
			return
		}

		dom.ConsoleLog("APPS", res.Apps)

		comp.state.Set("state:apps:changed", res.Apps)
	})
}

func (comp *AppsListComponent) updateApps() {

	raw := comp.state.Get("state:apps:changed")
	if raw == nil {
		return
	}

	apps, ok := raw.([]service.App)
	if !ok {
		dom.ConsoleError("appsListComponent changed error. cannot parse apps list")
		return
	}

	comp.apps = apps

	container := comp.GetContainer("appsList")
	if container.IsNull() {
		dom.ConsoleError("appsList component container is null", comp.ID())
		return
	}

	dom.SetInnerHTML(container, "")

	for _, app := range apps {
		appLauncher := NewAppLauncherComponent(comp.DI, AppComponentRegKey.WithExtra("-", app.ID), app)
		appLauncher.Mount(&container)
	}

	// All prepared! We can display the dashboard
	comp.state.Set("state:dashboard:ready", true)
}
