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
	element      js.Value
	fields       *LoginComponentFields
	unwatchState []func()
}

var _ component.Component = (*LoginComponent)(nil)
var _ component.EventControl = (*LoginComponent)(nil)
var _ component.StateControl = (*LoginComponent)(nil)

func NewLoginComponent(di *deps.Dependency) *LoginComponent {
	comp := &LoginComponent{
		di:     di,
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

func (comp *LoginComponent) init() {
	comp.BindEvents()
	comp.BindStateWatch()
}

func (comp *LoginComponent) stateErrored(_, newValue interface{}) {

	errContainer := dom.QuerySelector(`[data-dock="errorMessage"]`)
	errContainer.Set("innerHTML", "")
	dom.AddClass(errContainer, "hidden")

	errFields := newValue.(*ErrorManagerFields)
	if errFields == nil {
		return
	}

	errMgr := NewErrorManager(comp.di)
	errMgr.SetNamespace("auth/auth/errorMsg")
	errMgr.SetFields(errFields)

	if err := errMgr.Render(); err != nil {
		dom.ConsoleError("Cannot render errMsg", err.Error())
		return
	}

	if err := errMgr.Mount(&errContainer); err != nil {
		dom.ConsoleError("Cannot mount errMsg", err.Error())
		return
	}

	dom.RemoveClass(errContainer, "hidden")
}
