package components

import (
	"bytes"
	"errors"
	"orbital/web/wasm/pkg/component"
	"orbital/web/wasm/pkg/deps"
	"orbital/web/wasm/pkg/dom"
	"orbital/web/wasm/pkg/state"
	"syscall/js"
)

type TaskbarComponent struct {
	di           *deps.Dependency
	docks        map[string]js.Value
	element      js.Value
	state        *state.State
	unwatchState []func()
}

var _ component.ContainerComponent = (*TaskbarComponent)(nil)
var _ component.StateControl = (*TaskbarComponent)(nil)

func NewTaskbarComponent(di *deps.Dependency) *TaskbarComponent {
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

	return nil
}

func (comp *TaskbarComponent) Unmount() error {
	if !comp.element.IsNull() {
		dom.RemoveElement(comp.element)
		comp.element = js.Null()
	}

	comp.UnbindStateWatch()

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

func (comp *TaskbarComponent) BindStateWatch() {
	comp.state.Set("state:taskbar:currentMode", "")

	var unwatchFn func()

	unwatchFn = comp.state.Watch("state:taskbar:currentMode", comp.stateTaskbarCurrentMode)

	comp.unwatchState = append(comp.unwatchState, unwatchFn)
}

func (comp *TaskbarComponent) UnbindStateWatch() {
	for _, unwatchFn := range comp.unwatchState {
		unwatchFn()
	}
}

func (comp *TaskbarComponent) init() {
	comp.BindStateWatch()
}

func (comp *TaskbarComponent) stateTaskbarCurrentMode(oldLevel, newLevel interface{}) {
	newLvl := newLevel.(string)

	switch newLvl {
	case "guest":
		dom.ConsoleLog("Set taskbar in guest mode")

		container := comp.GetContainer("startMenu")
		if container.IsNull() {
			dom.ConsoleError("overlay component container is null", "login")
			return
		}
		dom.SetInnerHTML(container, "")

		startMenu := NewTaskbarStartComponent(comp.di)
		_ = startMenu.Render()
		_ = startMenu.Mount(&container)

	case "user":
		dom.ConsoleLog("Set taskbar in user mode")
	case "admin":
		dom.ConsoleLog("Set taskbar in admin mode")
	}
}
