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

type OverlayComponentFields struct {
	Title string
	Icon  string
}

func (fields *OverlayComponentFields) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"title": fields.Title,
		"icon":  fields.Icon,
	}
}

type OverlayComponent struct {
	di           *deps.Dependency
	docks        map[string]js.Value
	element      js.Value
	fields       OverlayComponentFields
	unwatchState []func()
	state        *state.State
}

var _ component.ContainerComponent = (*OverlayComponent)(nil)
var _ component.StateControl = (*OverlayComponent)(nil)

func NewOverlayComponent(di *deps.Dependency) *OverlayComponent {
	comp := &OverlayComponent{
		di:     di,
		fields: OverlayComponentFields{},
		docks:  make(map[string]js.Value),
		state:  di.State(),
	}

	comp.init()
	return comp
}

func (comp *OverlayComponent) ID() string {
	return "overlay"
}

func (comp *OverlayComponent) Namespace() string {
	return "orbital/overlay/overlay"
}

func (comp *OverlayComponent) Mount(container *js.Value) error {
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

func (comp *OverlayComponent) Unmount() error {
	if !comp.element.IsNull() {
		dom.RemoveElement(comp.element)
		comp.element = js.Null()
	}

	comp.unbindUIEvents()
	return nil
}

func (comp *OverlayComponent) Render() error {
	tpl, err := comp.di.TplRegistry().Get(comp.Namespace())
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, comp.fields.ToMap()); err != nil {
		return err
	}

	comp.element = dom.CreateElementFromString(buf.String())
	comp.SetContainers()

	return nil
}

func (comp *OverlayComponent) SetFields(fields OverlayComponentFields) {
	comp.fields = fields
}

func (comp *OverlayComponent) BindStateWatch() {
	comp.state.Set("state:overlay:activeComponent", "")

	var unwatchFn func()

	unwatchFn = comp.state.Watch("state:overlay:activeComponent", comp.stateOverlayActiveComponent)

	comp.unwatchState = append(comp.unwatchState, unwatchFn)
}

func (comp *OverlayComponent) UnbindStateWatch() {
	for _, unwatchFn := range comp.unwatchState {
		unwatchFn()
	}
}

func (comp *OverlayComponent) GetContainer(name string) js.Value {
	container, ok := comp.docks[name]
	if !ok {
		return js.Null()
	}

	if container.IsNull() {
		return js.Null()
	}

	return container
}

func (comp *OverlayComponent) SetContainers() {
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

func (comp *OverlayComponent) init() {
	comp.BindStateWatch()
}

func (comp *OverlayComponent) stateOverlayActiveComponent(oldActiveComp, newActiveComp interface{}) {
	newComp := newActiveComp.(string)
	oldComp := oldActiveComp.(string)

	dom.ConsoleLog("[OverlayComponent] state", oldComp, newComp)

	if newComp == oldComp {
		return
	}

	switch newActiveComp {
	case "login":

		comp.SetFields(OverlayComponentFields{
			Title: "Login",
			Icon:  "fa-lock",
		})

		container := comp.GetContainer("overlayData")
		if container.IsNull() {
			dom.ConsoleError("overlay component container is null", "login")
			return
		}
		dom.SetInnerHTML(container, "")

		loginComp := NewLoginComponent(comp.di)
		if err := loginComp.Render(); err != nil {
			dom.ConsoleError("overlay component login render error", err.Error())
			return
		}

		if err := loginComp.Mount(&container); err != nil {
			dom.ConsoleError("overlay component login mounting error", err.Error())
			return
		}

		dom.RemoveClass(comp.element, "hide")
	case "register":
		dom.ConsoleLog("Register not implemented")
	}
}
