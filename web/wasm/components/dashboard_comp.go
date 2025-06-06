package components

import (
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/dom"
	"orbital/web/wasm/pkg/state"
	"orbital/web/wasm/pkg/transport"
	"orbital/web/wasm/service"
)

const (
	DashboardComponentRegKey RegKey = "dashboardComponent"
)

type DashboardComponent struct {
	*BaseComponent
	state *state.State

	appsSvc *service.AppsService
}

func NewDashboardComponent(di *orbital.Dependency) *DashboardComponent {
	base := NewBaseComponent(di, DashboardComponentRegKey, "dashboard/main/default")

	appsSvc := orbital.MustGetService[*service.AppsService](di, service.AppsServiceKey)

	comp := &DashboardComponent{
		BaseComponent: base,
		state:         di.State,
		appsSvc:       appsSvc,
	}

	comp.OnMount(func() {
		comp.toggle("loading")

		comp.state.Set("state:dashboard:ready", true)

	})

	comp.state.Watch("state:dashboard:ready", func(oldV, newV interface{}) {
		comp.loadApps()
	})

	return comp
}

func (comp *DashboardComponent) ID() RegKey {
	return DashboardComponentRegKey
}

func (comp *DashboardComponent) loadApps() {
	var async transport.Async
	async.Async(func() {
		res, err := comp.appsSvc.List(service.ListReq{})
		if err != nil {
			dom.ConsoleError(err)
			return
		}

		dom.ConsoleLog("APP", res.Apps)
		comp.toggle("dashboard")
	})
}

func (comp *DashboardComponent) toggle(screen string) {
	loadingElem := dom.QuerySelector("[data-item='loading']")
	dashboardElem := dom.QuerySelector("[data-item='dashboard']")

	switch screen {
	case "loading":
		if !loadingElem.IsNull() {
			dom.RemoveClass(loadingElem, "hide")
		}
		if !dashboardElem.IsNull() {
			dom.AddClass(dashboardElem, "hide")
		}

	case "dashboard":
		if !loadingElem.IsNull() {
			dom.AddClass(loadingElem, "hide")
		}
		if !dashboardElem.IsNull() {
			dom.RemoveClass(dashboardElem, "hide")
		}
	}
}
