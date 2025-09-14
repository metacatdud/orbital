package components

import (
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/dom"
	"orbital/web/wasm/pkg/events"
	"orbital/web/wasm/pkg/state"
	"syscall/js"
)

const (
	MainComponentRegKey RegKey = "mainComponent"
)

type MainComponent struct {
	*BaseComponent
	di            *orbital.Dependency
	events        *events.Event
	state         *state.State
	dashboardComp Component
	taskbarComp   Component
	overlayComp   Component

	isDashboardMounted bool
}

func NewMainComponent(di *orbital.Dependency) *MainComponent {
	base := NewBaseComponent(di, MainComponentRegKey, "orbital/main/orbital")

	taskbarComp := NewTaskbarComponent(di)

	comp := &MainComponent{
		BaseComponent: base,
		di:            di,
		events:        di.Events,
		state:         di.State,
		taskbarComp:   taskbarComp,
	}

	comp.events.On("evt:overlay:show", func(data OverlayConfig) {
		comp.ShowOverlay(data)
	})

	comp.state.Watch("state:isAuthenticated", func(oldV, newV any) {
		comp.onAuthChanged(newV)
	})

	comp.OnMount(func() {
		comp.mountTaskbar()

		// Check if dashboard can be mounted during mount time
		isAuth := false
		if isAuthRaw := comp.state.Get("state:isAuthenticated"); isAuthRaw != nil {
			isAuth = isAuthRaw.(bool)
		}

		comp.mountDashboard(isAuth)
	})

	return comp
}

func (comp *MainComponent) ID() RegKey {
	return MainComponentRegKey
}

func (comp *MainComponent) Mount(container *js.Value) error {
	// TODO: Remove the node all together maybe?
	dom.QuerySelector("#loading-screen").
		Get("classList").
		Call("add", "hide")

	return comp.BaseComponent.Mount(container)
}

func (comp *MainComponent) Unmount() error {
	if comp.dashboardComp != nil {
		err := comp.dashboardComp.Unmount()
		if err != nil {
			return err
		}
	}
	if comp.taskbarComp != nil {
		err := comp.taskbarComp.Unmount()
		if err != nil {
			return err
		}
	}

	return comp.BaseComponent.Unmount()
}

func (comp *MainComponent) ShowOverlay(overlayConfig OverlayConfig) {
	container := comp.GetContainer("overlay")
	if container.IsNull() {
		dom.ConsoleError("[MainComponent] ShowOverlay. container is null")
		return
	}

	// Unmount any previous overlay if exists
	// This should not happen. Always make sure to call HideOverlay when
	// done using overlay.
	if comp.overlayComp != nil {
		comp.overlayComp.Unmount()
		comp.overlayComp = nil
	}

	dom.SetInnerHTML(container, "")

	comp.overlayComp = NewOverlayComponent(comp.di, overlayConfig)
	if err := comp.overlayComp.Mount(&container); err != nil {
		dom.ConsoleError("[MainComponent] ShowOverlay", err.Error())
		return
	}

	dom.RemoveClass(container, "hide")
}

func (comp *MainComponent) HideOverlay() {

	container := comp.GetContainer("overlay")
	if container.IsNull() {
		dom.ConsoleError("[MainComponent] ShowOverlay. container is null")
		return
	}

	dom.AddClass(container, "hide")
	if comp.overlayComp != nil {
		comp.overlayComp.Unmount()
		comp.overlayComp = nil
	}

	dom.SetInnerHTML(container, "")
}

func (comp *MainComponent) mountDashboard(shouldMount bool) {
	container := comp.GetContainer("dashboard")
	if container.IsNull() {
		dom.ConsoleError("[dashboard] mountDashboard. container is null")
		return
	}

	var err error
	if shouldMount && !comp.isDashboardMounted {
		comp.dashboardComp = NewDashboardComponent(comp.di)
		//if err != nil {
		//	dom.ConsoleError("[MainComponent] mountDashboard", err.Error())
		//	return
		//}

		if err = comp.dashboardComp.Mount(&container); err != nil {
			dom.ConsoleError("[MainComponent] mountDashboard", err.Error())
			return
		}

		comp.isDashboardMounted = true
		return
	}

	if !shouldMount && comp.isDashboardMounted {
		if comp.dashboardComp != nil {
			_ = comp.dashboardComp.Unmount()
		}
		comp.isDashboardMounted = false
	}
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

func (comp *MainComponent) onAuthChanged(val any) {
	isAuth := false
	if val != nil {
		isAuth = val.(bool)
	}

	// Cleanup everything first
	dashContainer := comp.GetContainer("dashboard")
	overlayContainer := comp.GetContainer("overlay")

	if !dashContainer.IsNull() {
		dom.SetInnerHTML(dashContainer, "")
	}

	if !overlayContainer.IsNull() {
		comp.HideOverlay()
	}

	comp.mountDashboard(isAuth)
	comp.setupAsUnauthenticated(isAuth)
}

// setupAsUnauthenticated will show an overlay in case the app has no session
func (comp *MainComponent) setupAsUnauthenticated(isAuth bool) {

	// If there is a session, skip!
	if isAuth {
		return
	}

	login := NewLoginComponent(comp.di)
	comp.ShowOverlay(OverlayConfig{
		Child:   login,
		Title:   login.Title(),
		Icon:    login.Icon(),
		Actions: []string{"close"},
		Css:     []string{"small"},
		OnClose: nil,
	})
	comp.isDashboardMounted = false
}
