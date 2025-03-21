package components

import (
	"bytes"
	"errors"
	"orbital/web/w2/orbital"
	"orbital/web/w2/pkg/dom"
	"orbital/web/w2/pkg/state"
	"syscall/js"
)

const (
	OverlayCompKey = "overlayComponent"
)

type OverlayComponentFields struct {
	Title    string
	Icon     string
	CompName string
}

func (fields *OverlayComponentFields) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"title":    fields.Title,
		"icon":     fields.Icon,
		"compName": fields.CompName,
	}
}

type OverlayComp struct {
	BaseComp
	di           *orbital.Dependency
	element      js.Value
	child        orbital.Component
	state        *state.State
	unwatchState []func()
}

// Implementation checklist
var _ orbital.ContainerComponent = (*OverlayComp)(nil)
var _ orbital.StateControl = (*OverlayComp)(nil)

func NewOverlayComp(di *orbital.Dependency) *OverlayComp {
	comp := &OverlayComp{
		BaseComp: BaseComp{docks: make(map[string]js.Value)},
		di:       di,
		state:    di.State(),
	}

	// Init
	comp.BindStateWatch()

	return comp
}

func (comp *OverlayComp) ID() string {
	return OverlayCompKey
}

func (comp *OverlayComp) Namespace() string {
	return "orbital/overlay/overlay"
}

func (comp *OverlayComp) Mount(container *js.Value) error {
	if !container.Truthy() {
		return errors.New("container does not exist")
	}

	if comp.element.IsNull() {
		return errors.New("element is missing")
	}

	dom.AppendChild(*container, comp.element)

	return nil
}

func (comp *OverlayComp) Unmount() error {
	if !comp.element.IsNull() {
		dom.RemoveElement(comp.element)
		comp.element = js.Null()
	}

	comp.UnbindStateWatch()
	return nil
}

func (comp *OverlayComp) Render() error {
	tpl, err := comp.di.TplRegistry().Get(comp.Namespace())
	if err != nil {
		return err
	}

	data := comp.state.Get("state:overlay:activeComponent")
	if data != nil {
		dom.ConsoleLog("overlay data", data)
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, nil); err != nil {
		return err
	}

	comp.element = dom.CreateElementFromString(buf.String())

	comp.SetContainers(comp.element)

	return nil
}

func (comp *OverlayComp) BindStateWatch() {
	comp.state.Set("state:overlay:activeComponent", OverlayComponentFields{})
	var unwatchFn func()

	unwatchFn = comp.state.Watch("state:overlay:activeComponent", comp.stateOverlayActiveComponent)
	comp.unwatchState = append(comp.unwatchState, unwatchFn)
}

func (comp *OverlayComp) UnbindStateWatch() {
	for _, unwatchFn := range comp.unwatchState {
		unwatchFn()
	}
}

func (comp *OverlayComp) stateOverlayActiveComponent(oldActiveComp, newActiveComp interface{}) {
	newComp := newActiveComp.(OverlayComponentFields)
	oldComp := oldActiveComp.(OverlayComponentFields)

	dom.ConsoleLog("Active Comp", newComp.ToMap(), oldComp.ToMap())
}
