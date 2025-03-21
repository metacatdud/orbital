package components

import (
	"bytes"
	"errors"
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/dom"
	"syscall/js"
)

type ErrorManagerFields struct {
	Type    string
	Message string
}

func (fields *ErrorManagerFields) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type":    fields.Type,
		"message": fields.Message,
	}
}

type ErrorManager struct {
	di      *orbital.Dependency
	fields  ErrorManagerFields
	element js.Value
}

var _ orbital.Component = (*ErrorManager)(nil)

func NewErrorManager(di *orbital.Dependency) *ErrorManager {
	comp := &ErrorManager{
		di:     di,
		fields: ErrorManagerFields{},
	}

	comp.init()

	return comp
}

func (comp *ErrorManager) Mount(container *js.Value) error {
	if !container.Truthy() {
		return errors.New("container does not exist")
	}

	if comp.element.IsNull() {
		return errors.New("element is missing")
	}

	dom.AppendChild(*container, comp.element)

	return nil
}

func (comp *ErrorManager) Unmount() error {
	if !comp.element.IsNull() {
		dom.RemoveElement(comp.element)
		comp.element = js.Null()
	}

	return nil
}

func (comp *ErrorManager) ID() string {
	return "error-manager"
}

func (comp *ErrorManager) Namespace() string {
	return "auth/auth/errorMsg"
}

func (comp *ErrorManager) Render() error {
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

func (comp *ErrorManager) SetFields(fields ErrorManagerFields) {
	comp.fields = fields
}

func (comp *ErrorManager) init() {}
