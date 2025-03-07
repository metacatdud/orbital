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

// Error manager component is a special component which can be used
// in a variety of ways:
// - Attach an error in a specific placeholder
// - Attach in a specific field such as `data-error-for="email"`
//
// This error must be used as subcomponent for other components

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
	di        *deps.Dependency
	fields    *ErrorManagerFields
	namespace string
	element   js.Value
}

var _ component.Component = (*ErrorManager)(nil)

func NewErrorManager(di *deps.Dependency) *ErrorManager {
	comp := &ErrorManager{
		di:     di,
		fields: &ErrorManagerFields{},
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
	return comp.namespace
}

func (comp *ErrorManager) Render() error {
	tpl, err := comp.di.TplRegistry().Get(comp.Namespace())
	if err != nil {
		return err
	}

	// hackish merge to check proper state and fields merging
	mergedData := state.MergeStateWithData(
		comp.di.State().GetAll(),
		comp.fields.ToMap(),
	)

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, mergedData); err != nil {
		return err
	}

	comp.element = dom.CreateElementFromString(buf.String())

	return nil
}

func (comp *ErrorManager) SetFields(fields *ErrorManagerFields) {
	comp.fields = fields
}

func (comp *ErrorManager) SetNamespace(namespace string) {
	comp.namespace = namespace
}

func (comp *ErrorManager) init() {}
