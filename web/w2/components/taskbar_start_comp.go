package components

import (
	"bytes"
	"errors"
	"orbital/web/w2/orbital"
	"orbital/web/w2/pkg/dom"
	"syscall/js"
)

const (
	TaskbarStartCompKey = "taskbarStartComponent"
)

type TaskbarStartComp struct {
	BaseComp
	di      *orbital.Dependency
	element js.Value
}

var _ orbital.ContainerComponent = (*TaskbarStartComp)(nil)

func NewTaskbarStartComp(di *orbital.Dependency) *TaskbarStartComp {
	comp := &TaskbarStartComp{
		BaseComp: BaseComp{docks: make(map[string]js.Value)},
		di:       di,
	}

	return comp
}

func (comp *TaskbarStartComp) ID() string {
	return TaskbarStartCompKey
}

func (comp *TaskbarStartComp) Namespace() string {
	return "orbital/taskbar/taskbarStartContent"
}

func (comp *TaskbarStartComp) Mount(container *js.Value) error {
	if !container.Truthy() {
		return errors.New("container does not exist")
	}

	if comp.element.IsNull() {
		return errors.New("element is missing")
	}

	dom.AppendChild(*container, comp.element)

	return nil
}

func (comp *TaskbarStartComp) Unmount() error {
	if !comp.element.IsNull() {
		dom.RemoveElement(comp.element)
		comp.element = js.Null()
	}

	return nil
}

func (comp *TaskbarStartComp) Render() error {
	tpl, err := comp.di.TplRegistry().Get(comp.Namespace())
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, nil); err != nil {
		return err
	}

	comp.element = dom.CreateElementFromString(buf.String())

	comp.SetContainers(comp.element)

	return nil
}
