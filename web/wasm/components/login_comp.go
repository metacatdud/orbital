package components

import (
	"bytes"
	"errors"
	"orbital/web/wasm/pkg/component"
	"orbital/web/wasm/pkg/deps"
	"orbital/web/wasm/pkg/dom"
	"syscall/js"
)

type LoginComponentFields struct {
	PrivateKey string
}

func (field *LoginComponentFields) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"privateKey": field.PrivateKey,
	}
}

type LoginComponent struct {
	di           *deps.Dependency
	docks        map[string]js.Value
	element      js.Value
	fields       *LoginComponentFields
	unwatchState []func()
}

var _ component.ContainerComponent = (*LoginComponent)(nil)
var _ component.EventControl = (*LoginComponent)(nil)
var _ component.StateControl = (*LoginComponent)(nil)

func NewLoginComponent(di *deps.Dependency) *LoginComponent {
	comp := &LoginComponent{
		di:     di,
		docks:  make(map[string]js.Value),
		fields: &LoginComponentFields{},
	}

	comp.init()

	return comp

}

func (comp *LoginComponent) ID() string {
	return "login"
}

func (comp *LoginComponent) Namespace() string {
	return "auth/auth/default"
}

func (comp *LoginComponent) Mount(container *js.Value) error {
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

func (comp *LoginComponent) Unmount() error {
	if !comp.element.IsNull() {
		dom.RemoveElement(comp.element)
		comp.element = js.Null()
	}

	comp.unbindUIEvents()
	comp.UnbindStateWatch()

	return nil
}

func (comp *LoginComponent) Render() error {
	dom.ConsoleLog("- Rendering", comp.ID())

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

func (comp *LoginComponent) BindStateWatch() {
	unwatchAuthFn := comp.di.State().Watch("state:auth:errored", comp.stateErrored)

	comp.unwatchState = append(comp.unwatchState, unwatchAuthFn)
}

func (comp *LoginComponent) UnbindStateWatch() {
	for _, unwatchFn := range comp.unwatchState {
		unwatchFn()
	}
}

func (comp *LoginComponent) GetContainer(name string) js.Value {
	container, ok := comp.docks[name]
	if !ok {
		return js.Null()
	}

	if container.IsNull() {
		return js.Null()
	}

	return container
}

func (comp *LoginComponent) SetContainers() {
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

func (comp *LoginComponent) init() {
	comp.BindEvents()
	comp.BindStateWatch()
}

func (comp *LoginComponent) stateErrored(_, newValue interface{}) {

	container := comp.GetContainer("errorMessage")
	if container.IsNull() {
		dom.ConsoleError("overlay component container is null", "errorComp")
		return
	}
	dom.AddClass(container, "hidden")
	dom.SetInnerHTML(container, "")

	errFields := newValue.(ErrorManagerFields)

	errMgr := NewErrorManager(comp.di)
	errMgr.SetFields(errFields)

	if err := errMgr.Render(); err != nil {
		dom.ConsoleError("Cannot render errMsg", err.Error())
		return
	}

	if err := errMgr.Mount(&container); err != nil {
		dom.ConsoleError("Cannot mount errMsg", err.Error())
		return
	}

	dom.RemoveClass(container, "hidden")
}
