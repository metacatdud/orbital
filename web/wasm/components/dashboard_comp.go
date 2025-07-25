package components

import (
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/dom"
	"orbital/web/wasm/pkg/state"
)

const (
	DashboardComponentRegKey RegKey = "dashboardComponent"
)

type DashboardComponent struct {
	*BaseComponent
	state   *state.State
	appList Component
}

func NewDashboardComponent(di *orbital.Dependency) *DashboardComponent {
	base := NewBaseComponent(di, DashboardComponentRegKey, "dashboard/main/default")

	// Child components
	appList := NewAppsListComponent(di)

	comp := &DashboardComponent{
		BaseComponent: base,
		state:         di.State,
		appList:       appList,
	}

	//comp.state.Set("state:dashboard:ready", false)

	comp.OnMount(func() {
		comp.toggle("loading")
		comp.mountAppList()
	})

	comp.state.Watch("state:dashboard:ready", func(oldV, newV interface{}) {
		if newV.(bool) {
			comp.toggle("dashboard")
		}
	})

	return comp
}

func (comp *DashboardComponent) ID() RegKey {
	return DashboardComponentRegKey
}

//func (comp *DashboardComponent) Mount(container *js.Value) error {
//	if !container.Truthy() {
//		return fmt.Errorf("dashboard component does not mount")
//	}
//
//	if err := comp.BaseComponent.Mount(container); err != nil {
//		return err
//	}
//
//	comp.mountAppList()
//	return nil
//}

func (comp *DashboardComponent) Unmount() error {
	if comp.appList != nil {
		comp.appList.Unmount()
	}

	return comp.BaseComponent.Unmount()
}

func (comp *DashboardComponent) mountAppList() {
	container := comp.GetContainer("appList")
	if container.IsNull() {
		dom.ConsoleError("dashboard component container is null", comp.ID())
		return
	}

	dom.SetInnerHTML(container, "")

	if err := comp.appList.Mount(&container); err != nil {
		dom.ConsoleError("dashboard component cannot be mounted", err.Error(), comp.ID())
		return
	}
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
