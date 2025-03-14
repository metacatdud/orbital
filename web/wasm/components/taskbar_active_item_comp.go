package components

import (
	"bytes"
	"errors"
	"orbital/web/wasm/pkg/component"
	"orbital/web/wasm/pkg/deps"
	"orbital/web/wasm/pkg/dom"
	"syscall/js"
)

type TaskbarActiveItemFields struct {
	Id    string
	Title string
	Icon  string
}

func (fields *TaskbarActiveItemFields) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":    fields.Id,
		"title": fields.Title,
		"icon":  fields.Icon,
	}
}

type TaskbarActiveItemComponent struct {
	di      *deps.Dependency
	fields  TaskbarActiveItemFields
	element js.Value
}

var _ component.Component = (*TaskbarActiveItemComponent)(nil)

func NewTaskbarActiveItemComponent(di *deps.Dependency) *TaskbarActiveItemComponent {
	comp := &TaskbarActiveItemComponent{
		di:     di,
		fields: TaskbarActiveItemFields{},
	}

	comp.init()

	return comp
}

func (comp *TaskbarActiveItemComponent) ID() string {
	return "taskbarActiveItem"
}

func (comp *TaskbarActiveItemComponent) Namespace() string {
	return "orbital/taskbar/taskbarActiveAppsItem"
}

func (comp *TaskbarActiveItemComponent) Mount(container *js.Value) error {
	if !container.Truthy() {
		return errors.New("container does not exist")
	}

	if comp.element.IsNull() {
		return errors.New("element is missing")
	}

	dom.AppendChild(*container, comp.element)

	return nil
}

func (comp *TaskbarActiveItemComponent) Unmount() error {
	if !comp.element.IsNull() {
		dom.RemoveElement(comp.element)
		comp.element = js.Null()
	}

	return nil
}

func (comp *TaskbarActiveItemComponent) Render() error {
	tpl, err := comp.di.TplRegistry().Get(comp.Namespace())
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, comp.fields.ToMap()); err != nil {
		return err
	}

	comp.element = dom.CreateElementFromString(buf.String())

	return nil
}

func (comp *TaskbarActiveItemComponent) SetFields(fields TaskbarActiveItemFields) {
	comp.fields = fields
}

func (comp *TaskbarActiveItemComponent) init() {}
