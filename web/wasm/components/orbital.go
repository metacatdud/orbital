package components

import (
	"errors"
	"orbital/web/wasm/pkg/component"
	"orbital/web/wasm/pkg/deps"
	"orbital/web/wasm/pkg/dom"
	"syscall/js"
)

type OrbitalComponent struct {
	di           *deps.Dependency
	element      *js.Value
	unwatchState []func()
}

// Implementation checklist
var _ component.Component = (*OrbitalComponent)(nil)
var _ component.StateControl = (*OrbitalComponent)(nil)

func NewOrbitalComponent(di *deps.Dependency) *OrbitalComponent {
	o := &OrbitalComponent{
		di: di,
	}

	o.init()

	return o
}

func (comp *OrbitalComponent) ID() string {
	return ""
}

func (comp *OrbitalComponent) Namespace() string {
	return ""
}

func (comp *OrbitalComponent) Render() error {
	return nil
}

func (comp *OrbitalComponent) Mount(container *js.Value) error {
	if !container.Truthy() {
		return errors.New("container does not exist")
	}

	loadingElem := dom.QuerySelector("#loading")
	if !loadingElem.IsNull() {
		dom.RemoveElement(loadingElem)
	}

	// For this component only the container already exist as the #root div from index.html
	comp.element = container

	return nil
}

func (comp *OrbitalComponent) Unmount() error {
	comp.UnbindStateWatch()
	return nil
}

func (comp *OrbitalComponent) BindStateWatch() {
	unwatchAuthFn := comp.di.State().Watch("state:orbital:authenticated", func(oldValue, newValue interface{}) {
		if newValue.(bool) {
			comp.renderDashboard()
			return
		}

		comp.renderLogin()
	})

	comp.unwatchState = append(comp.unwatchState, unwatchAuthFn)
}

func (comp *OrbitalComponent) UnbindStateWatch() {
	for _, unwatchFn := range comp.unwatchState {
		unwatchFn()
	}
}

func (comp *OrbitalComponent) init() {
	comp.BindStateWatch()
}

func (comp *OrbitalComponent) renderLogin() {
	if comp.element == nil || comp.element.IsUndefined() {
		dom.ConsoleError("[OrbitalComponent] not mounted properly")
		return
	}

	comp.element.Set("innerHTML", "")

	loginComponent := NewLoginComponent(comp.di)

	if err := loginComponent.Render(); err != nil {
		dom.ConsoleError("[OrbitalComponent] Cannot render LoginComponent", err.Error())
		return
	}

	if err := loginComponent.Mount(comp.element); err != nil {
		dom.ConsoleError("[OrbitalComponent] Cannot mount LoginComponent", err.Error())
	}
}

func (comp *OrbitalComponent) renderDashboard() {
	if comp.element == nil || comp.element.IsUndefined() {
		dom.ConsoleError("[OrbitalComponent] not mounted properly")
		return
	}

	comp.element.Set("innerHTML", "")
	dashComponent := NewDashboardComponent(comp.di)

	if err := dashComponent.Render(); err != nil {
		dom.ConsoleError("[OrbitalComponent] Cannot render DashboardComponent", err.Error())
	}

	if err := dashComponent.Mount(comp.element); err != nil {
		dom.ConsoleError("[OrbitalComponent] Cannot mount DashboardComponent", err.Error())
	}
}

func (comp *OrbitalComponent) renderRegister() {
	// TODO: Implement renderRegister
}
