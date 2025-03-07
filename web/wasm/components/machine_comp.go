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

type MachineComponentFields struct{}

func (field *MachineComponentFields) ToMap() map[string]interface{} {
	return map[string]interface{}{}
}

type MachineComponent struct {
	di           *deps.Dependency
	element      js.Value
	fields       MachineComponentFields
	unwatchState []func()
}

var _ component.Component = (*MachineComponent)(nil)

func NewMachineComponent(di *deps.Dependency) *MachineComponent {
	comp := &MachineComponent{
		di:     di,
		fields: MachineComponentFields{},
	}

	comp.init()
	return comp
}

func (comp *MachineComponent) ID() string {
	return "machine"
}

func (comp *MachineComponent) Namespace() string {
	return "dashboard/machines/machinesList"
}

func (comp *MachineComponent) Mount(container *js.Value) error {
	dom.ConsoleLog("- Mounting", comp.ID())

	if !container.Truthy() {
		return errors.New("container does not exist")
	}

	if comp.element.IsNull() {
		return errors.New("element is missing")
	}

	dom.AppendChild(*container, comp.element)

	// Add bindUIEvents here

	return nil
}

func (comp *MachineComponent) Unmount() error {
	if !comp.element.IsNull() {
		dom.RemoveElement(comp.element)
		comp.element = js.Null()
	}

	// Add unbindUIEvents here
	return nil
}

func (comp *MachineComponent) Render() error {
	tpl, err := comp.di.TplRegistry().Get(comp.Namespace())
	if err != nil {
		return err
	}

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

func (comp *MachineComponent) SetFields(fields MachineComponentFields) {
	comp.fields = fields
}

func (comp *MachineComponent) init() {}
