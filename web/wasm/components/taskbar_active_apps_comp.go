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

type TaskbarActiveAppsComponent struct {
	di           *deps.Dependency
	docks        map[string]js.Value
	element      js.Value
	state        *state.State
	unwatchState []func()
}

var _ component.ContainerComponent = (*TaskbarActiveAppsComponent)(nil)

func NewTaskbarActiveAppsComponent(di *deps.Dependency) *TaskbarActiveAppsComponent {
	comp := &TaskbarActiveAppsComponent{
		di:    di,
		docks: make(map[string]js.Value),
		state: di.State(),
	}

	comp.init()

	return comp
}
func (comp *TaskbarActiveAppsComponent) ID() string {
	return "taskbarActiveApps"
}

func (comp *TaskbarActiveAppsComponent) Namespace() string {
	return "orbital/taskbar/taskbarActiveApps"
}

func (comp *TaskbarActiveAppsComponent) Mount(container *js.Value) error {
	if !container.Truthy() {
		return errors.New("container does not exist")
	}

	if comp.element.IsNull() {
		return errors.New("element is missing")
	}

	dom.AppendChild(*container, comp.element)

	return nil
}

func (comp *TaskbarActiveAppsComponent) Unmount() error {
	if !comp.element.IsNull() {
		dom.RemoveElement(comp.element)
		comp.element = js.Null()
	}

	//comp.UnbindStateWatch()

	return nil
}

func (comp *TaskbarActiveAppsComponent) Render() error {
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

func (comp *TaskbarActiveAppsComponent) GetContainer(name string) js.Value {
	container, ok := comp.docks[name]
	if !ok {
		return js.Null()
	}

	if container.IsNull() {
		return js.Null()
	}

	return container
}

func (comp *TaskbarActiveAppsComponent) SetContainers() {
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

func (comp *TaskbarActiveAppsComponent) init() {}
