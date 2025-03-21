package components

import (
	"bytes"
	"errors"
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/dom"
	"orbital/web/wasm/pkg/state"
	"syscall/js"
)

type TaskbarComponent struct {
	di           *orbital.Dependency
	docks        map[string]js.Value
	element      js.Value
	state        *state.State
	unwatchState []func()
}

var _ orbital.ContainerComponent = (*TaskbarComponent)(nil)

func NewTaskbarComponent(di *orbital.Dependency) *TaskbarComponent {
	comp := &TaskbarComponent{
		di:    di,
		docks: make(map[string]js.Value),
		state: di.State(),
	}

	comp.init()

	return comp
}

func (comp *TaskbarComponent) ID() string {
	return "taskbar"
}

func (comp *TaskbarComponent) Namespace() string {
	return "orbital/taskbar/taskbar"
}

func (comp *TaskbarComponent) Mount(container *js.Value) error {
	if !container.Truthy() {
		return errors.New("container does not exist")
	}

	if comp.element.IsNull() {
		return errors.New("element is missing")
	}

	dom.AppendChild(*container, comp.element)
	comp.bindUIEvents()

	return nil
}

func (comp *TaskbarComponent) Unmount() error {
	if !comp.element.IsNull() {
		dom.RemoveElement(comp.element)
		comp.element = js.Null()
	}

	comp.unbindUIEvents()

	return nil
}

func (comp *TaskbarComponent) Render() error {
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
	comp.renderChildComponents()

	return nil
}

func (comp *TaskbarComponent) GetContainer(name string) js.Value {
	container, ok := comp.docks[name]
	if !ok {
		return js.Null()
	}

	if container.IsNull() {
		return js.Null()
	}

	return container
}

func (comp *TaskbarComponent) SetContainers() {
	if comp.element.IsNull() {
		dom.ConsoleError("element is missing", comp.ID())
		return
	}

	dockingAreas := dom.QuerySelectorAllFromElement(comp.element, `[data-dock]`)
	for _, area := range dockingAreas {
		areaName := area.Get("dataset").Get("dock").String()
		comp.docks[areaName] = area
	}
}

func (comp *TaskbarComponent) init() {}

func (comp *TaskbarComponent) renderChildComponents() {
	comp.renderStartMenu()
	comp.renderActiveApps()
	comp.renderSystemTray()
}

func (comp *TaskbarComponent) renderStartMenu() {
	container := comp.GetContainer("startMenu")
	if container.IsNull() {
		dom.ConsoleError("overlay component container is null", "login")
		return
	}
	dom.SetInnerHTML(container, "")

	startMenu := NewTaskbarStartComponent(comp.di)
	_ = startMenu.Render()
	_ = startMenu.Mount(&container)
}

func (comp *TaskbarComponent) renderActiveApps() {
	container := comp.GetContainer("activeApps")
	if container.IsNull() {
		dom.ConsoleError("overlay component container is null", "login")
		return
	}
	dom.SetInnerHTML(container, "")
	activeApps := NewTaskbarActiveAppsComponent(comp.di)
	_ = activeApps.Render()
	_ = activeApps.Mount(&container)
}

func (comp *TaskbarComponent) renderSystemTray() {}
