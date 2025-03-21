package components

import (
	"bytes"
	"errors"
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/dom"
	"orbital/web/wasm/pkg/state"
	"syscall/js"
)

type OrbitalComponent struct {
	di           *orbital.Dependency
	docks        map[string]js.Value
	element      js.Value
	factory      *orbital.Factory
	unwatchState []func()
	state        *state.State
}

// Implementation checklist
var _ orbital.ContainerComponent = (*OrbitalComponent)(nil)
var _ orbital.StateControl = (*OrbitalComponent)(nil)

func NewOrbitalComponent(di *orbital.Dependency) *OrbitalComponent {
	o := &OrbitalComponent{
		di:      di,
		docks:   make(map[string]js.Value),
		factory: di.Factory(),
		state:   di.State(),
	}

	o.init()

	return o
}

func (comp *OrbitalComponent) ID() string {
	return "orbital"
}

func (comp *OrbitalComponent) Namespace() string {
	return "orbital/main/orbital"
}

func (comp *OrbitalComponent) Render() error {
	tpl, err := comp.di.TplRegistry().Get(comp.Namespace())
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, nil); err != nil {
		return err
	}

	comp.element = dom.CreateElementFromString(buf.String())

	comp.SetContainers()

	return nil
}

func (comp *OrbitalComponent) Mount(container *js.Value) error {
	if !container.Truthy() {
		return errors.New("container does not exist")
	}

	loadingElem := dom.QuerySelector("#loading")
	if !loadingElem.IsNull() {
		dom.RemoveElement(loadingElem)
	}

	if comp.element.IsNull() {
		return errors.New("element is missing")
	}

	dom.AppendChild(*container, comp.element)
	comp.mountChildComponents()

	return nil
}

func (comp *OrbitalComponent) Unmount() error {
	comp.UnbindStateWatch()
	return nil
}

func (comp *OrbitalComponent) BindStateWatch() {
	unwatchAuthFn := comp.state.Watch("state:orbital:authenticated", func(oldValue, newValue interface{}) {
		newAuthVal := newValue.(bool)
		if newAuthVal {

		}
	})

	comp.unwatchState = append(comp.unwatchState, unwatchAuthFn)
}

func (comp *OrbitalComponent) UnbindStateWatch() {
	for _, unwatchFn := range comp.unwatchState {
		unwatchFn()
	}
}

// The following two methods are helpers for docking various components into
// docking areas of the components it servers

func (comp *OrbitalComponent) GetContainer(name string) js.Value {
	container, ok := comp.docks[name]
	if !ok {
		return js.Null()
	}

	if container.IsNull() {
		return js.Null()
	}

	return container
}

func (comp *OrbitalComponent) SetContainers() {
	if comp.element.IsNull() {
		dom.ConsoleError("element is missing", comp.ID())
		return
	}

	// Set docking points
	// TODO: Improve validation and make sure these are not null
	dockingAreas := dom.QuerySelectorAllFromElement(comp.element, `[data-dock]`)
	if len(dockingAreas) == 0 {
		return
	}

	for _, area := range dockingAreas {
		areaName := area.Get("dataset").Get("dock").String()
		comp.docks[areaName] = area
	}
}

func (comp *OrbitalComponent) init() {
	comp.BindStateWatch()
}

func (comp *OrbitalComponent) mountChildComponents() {
	comp.mountDesktop()
	comp.mountOverlay()
	comp.mountTaskbar()
}

func (comp *OrbitalComponent) mountDesktop() {
	dom.ConsoleLog("[orbital] rendering desktop. Not implemented")
}

func (comp *OrbitalComponent) mountOverlay() {

	container := comp.GetContainer("overlay")
	if container.IsNull() {
		dom.ConsoleError("overlay component container is null", "login")
		return
	}
	dom.SetInnerHTML(container, "")

	overlayComponent := NewOverlayComponent(comp.di)
	_ = overlayComponent.Render()
	_ = overlayComponent.Mount(&container)
}

func (comp *OrbitalComponent) mountTaskbar() {
	container := comp.GetContainer("taskbar")
	if container.IsNull() {
		dom.ConsoleError("overlay component container is null", "login")
		return
	}
	dom.SetInnerHTML(container, "")

	taskbarComponent := NewTaskbarComponent(comp.di)
	_ = taskbarComponent.Render()
	_ = taskbarComponent.Mount(&container)
}
