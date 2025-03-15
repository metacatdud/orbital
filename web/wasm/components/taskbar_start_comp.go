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

type TaskbarStartComponent struct {
	di           *deps.Dependency
	docks        map[string]js.Value
	element      js.Value
	state        *state.State
	unwatchState []func()
}

var _ component.ContainerComponent = (*TaskbarStartComponent)(nil)
var _ component.StateControl = (*TaskbarStartComponent)(nil)

func NewTaskbarStartComponent(di *deps.Dependency) *TaskbarStartComponent {
	comp := &TaskbarStartComponent{
		di:    di,
		docks: make(map[string]js.Value),
		state: di.State(),
	}

	comp.init()

	return comp
}

func (comp *TaskbarStartComponent) ID() string {
	return "taskbarStartComponent"
}

func (comp *TaskbarStartComponent) Namespace() string {
	return "orbital/taskbar/taskbarStartContent"
}

func (comp *TaskbarStartComponent) Mount(container *js.Value) error {
	if !container.Truthy() {
		return errors.New("container does not exist")
	}

	if comp.element.IsNull() {
		return errors.New("element is missing")
	}

	dom.AppendChild(*container, comp.element)

	return nil
}

func (comp *TaskbarStartComponent) Unmount() error {
	if !comp.element.IsNull() {
		dom.RemoveElement(comp.element)
		comp.element = js.Null()
	}

	return nil
}

func (comp *TaskbarStartComponent) Render() error {
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

func (comp *TaskbarStartComponent) GetContainer(name string) js.Value {
	container, ok := comp.docks[name]
	if !ok {
		return js.Null()
	}

	if container.IsNull() {
		return js.Null()
	}

	return container
}

func (comp *TaskbarStartComponent) SetContainers() {
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

func (comp *TaskbarStartComponent) BindStateWatch() {
	comp.state.Set("state:taskbar:currentMode", "")

	var unwatchFn func()

	unwatchFn = comp.state.Watch("state:taskbar:currentMode", comp.stateTaskbarCurrentMode)

	comp.unwatchState = append(comp.unwatchState, unwatchFn)
}

func (comp *TaskbarStartComponent) UnbindStateWatch() {
	for _, unwatchFn := range comp.unwatchState {
		unwatchFn()
	}
}

func (comp *TaskbarStartComponent) init() {
	comp.BindStateWatch()
}

// stateTaskbarCurrentMode captures status change based on the role of current user
func (comp *TaskbarStartComponent) stateTaskbarCurrentMode(oldLevel, newLevel interface{}) {
	newLvl := newLevel.(string)

	switch newLvl {
	case "guest":
		dom.ConsoleLog("Set taskbar in guest mode")
	default:
		dom.ConsoleLog("Set taskbar validate against BE")
	}
}
