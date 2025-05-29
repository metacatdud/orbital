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
	overlayComp Component
	taskbarComp Component
}

func NewMainComponent(di *orbital.Dependency) *MainComponent {
	base := NewBaseComponent(di, MainComponentRegKey, "orbital/main/orbital")

	// TODO: Register children components here: taskbar, overlay, desktop
	overlayComp := NewOverlayComponent(di)
	taskbarComp := NewTaskbarComponent(di)

	comp := &MainComponent{
		BaseComponent: base,
		state:         di.State,
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
	comp.overlayComp.Unmount()
	comp.taskbarComp.Unmount()

	return comp.BaseComponent.Unmount()
}

func (comp *MainComponent) onMountHandler() {

	// Hide loading screen
	loadingElem := dom.QuerySelector("#loading")
	if !loadingElem.IsNull() {
		dom.RemoveElement(loadingElem)
	}

	comp.state.Watch("state:overlay:toggle", func(oldV, newV interface{}) {
		comp.toggleOverlay(newV.(bool))
	})

	// TODO: Add child components here
	comp.mountOverlay()
	comp.mountTaskbar()

}

func (comp *MainComponent) mountOverlay() {
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
