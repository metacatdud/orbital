package components

import (
	"fmt"
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/dom"
	"orbital/web/wasm/pkg/state"
	"syscall/js"
)

const (
	MainComponentRegKey RegKey = "mainComponent"
)

type MainComponent struct {
	*BaseComponent

	state *state.State

	//TODO: add components here
	dashboardComp Component
	overlayComp   Component
	taskbarComp   Component
}

func NewMainComponent(di *orbital.Dependency) *MainComponent {
	base := NewBaseComponent(di, MainComponentRegKey, "orbital/main/orbital")

	// TODO: Register children components here: taskbar, overlay, desktop
	dashboardComp := NewDashboardComponent(di)
	overlayComp := NewOverlayComponent(di)
	taskbarComp := NewTaskbarComponent(di)

	comp := &MainComponent{
		BaseComponent: base,
		state:         di.State,
		dashboardComp: dashboardComp,
		overlayComp:   overlayComp,
		taskbarComp:   taskbarComp,
	}

	// OnMount hook
	base.OnMount(comp.onMountHandler)

	return comp
}

func (comp *MainComponent) ID() RegKey {
	return MainComponentRegKey
}

func (comp *MainComponent) Mount(container *js.Value) error {
	if !container.Truthy() {
		return fmt.Errorf("orbital main component does not mount")
	}

	return comp.BaseComponent.Mount(container)
}

func (comp *MainComponent) Unmount() error {
	comp.dashboardComp.Unmount()
	comp.overlayComp.Unmount()
	comp.taskbarComp.Unmount()

	return comp.BaseComponent.Unmount()
}

func (comp *MainComponent) onMountHandler() {

	// Hide loading screen
	loadingElem := dom.QuerySelector("#loading-screen")
	if !loadingElem.IsNull() {
		dom.RemoveElement(loadingElem)
	}

	comp.state.Watch("state:overlay:toggle", func(oldV, newV interface{}) {
		comp.toggleOverlay(newV.(bool))
	})

	comp.state.Watch("state:isAuthenticated", func(oldV, newV interface{}) {
		if newV.(bool) {
			comp.overlayComp.Unmount()
			comp.mountDashboard(true)
		}
	})

	isAuth := false
	isAuthRaw := comp.state.Get("state:isAuthenticated")
	if isAuthRaw != nil {
		isAuth = isAuthRaw.(bool)
	}

	comp.mountDashboard(isAuth)
	comp.mountOverlay(isAuth)
	comp.mountTaskbar()
}

func (comp *MainComponent) mountDashboard(shouldInit bool) {
	if !shouldInit {
		return
	}

	container := comp.GetContainer("dashboard")
	if container.IsNull() {
		dom.ConsoleError("overlay component container is null", comp.ID())
		return
	}

	dom.SetInnerHTML(container, "")

	if err := comp.dashboardComp.Mount(&container); err != nil {
		dom.ConsoleError("overlay component cannot be mounted", err.Error(), comp.ID())
		return
	}
}

func (comp *MainComponent) toggleDashboard(show bool) {
	container := comp.GetContainer("dashboard")
	if container.IsNull() {
		dom.ConsoleError("dashboard component container is null", comp.ID())
		return
	}

	if show {
		dom.RemoveClass(container, "hide")
		return
	}

	dom.AddClass(container, "hide")
}

func (comp *MainComponent) mountOverlay(shouldInit bool) {
	if shouldInit {
		return
	}

	comp.state.Set("state:overlay:currentChild", LoginComponentRegKey)

	container := comp.GetContainer("overlay")
	if container.IsNull() {
		dom.ConsoleError("overlay component container is null", comp.ID())
		return
	}

	dom.SetInnerHTML(container, "")

	if err := comp.overlayComp.Mount(&container); err != nil {
		dom.ConsoleError("overlay component cannot be mounted", err.Error(), comp.ID())
		return
	}
}

func (comp *MainComponent) toggleOverlay(show bool) {
	container := comp.GetContainer("overlay")
	if container.IsNull() {
		dom.ConsoleError("overlay component container is null", comp.ID())
		return
	}

	if show {
		dom.RemoveClass(container, "hide")
		return
	}

	dom.AddClass(container, "hide")
}

func (comp *MainComponent) mountTaskbar() {
	container := comp.GetContainer("taskbar")
	if container.IsNull() {
		dom.ConsoleError("taskbar component container is null", comp.ID())
		return
	}

	dom.SetInnerHTML(container, "")

	if err := comp.taskbarComp.Mount(&container); err != nil {
		dom.ConsoleError("taskbar component cannot be mounted", err.Error(), comp.ID())
		return
	}
}
