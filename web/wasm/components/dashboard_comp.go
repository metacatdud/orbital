package components

import (
	"bytes"
	"errors"
	"orbital/orbital"
	"orbital/web/wasm/domain"
	"orbital/web/wasm/pkg/component"
	"orbital/web/wasm/pkg/deps"
	"orbital/web/wasm/pkg/dom"
	"orbital/web/wasm/pkg/events"
	"syscall/js"
)

type DashboardComponent struct {
	di               *deps.Dependency
	events           *events.Event
	element          js.Value
	unsubscribeEvent []func()
}

var _ component.Component = (*DashboardComponent)(nil)
var _ component.EventControl = (*DashboardComponent)(nil)

func NewDashboardComponent(di *deps.Dependency) *DashboardComponent {
	comp := &DashboardComponent{
		di:     di,
		events: di.Events(),
	}

	comp.init()

	return comp
}

func (comp *DashboardComponent) ID() string {
	return "dashboard"
}

func (comp *DashboardComponent) Namespace() string {
	return "dashboard/main/default"
}

func (comp *DashboardComponent) Mount(container *js.Value) error {
	dom.ConsoleLog("- Mounting", comp.ID())

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

func (comp *DashboardComponent) Unmount() error {
	if !comp.element.IsNull() {
		dom.RemoveElement(comp.element)
		comp.element = js.Null()
	}

	comp.unbindUIEvents()
	comp.UnbindEvents()

	return nil
}

func (comp *DashboardComponent) Render() error {
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

func (comp *DashboardComponent) BindEvents() {
	unsubEvMaUpd := comp.events.On("evt:machines:update", comp.eventMachine)
	comp.unsubscribeEvent = append(comp.unsubscribeEvent, unsubEvMaUpd)

	unsubEvMaErr := comp.events.On("evt:machines:error", comp.eventMachineError)
	comp.unsubscribeEvent = append(comp.unsubscribeEvent, unsubEvMaErr)
}

func (comp *DashboardComponent) UnbindEvents() {
	for _, unsubEv := range comp.unsubscribeEvent {
		unsubEv()
	}
}

func (comp *DashboardComponent) init() {
	comp.BindEvents()
}

func (comp *DashboardComponent) bindUIEvents() {
	dom.AddEventListener(`[data-action='toggleProfile']`, "click", comp.uiEventToggleProfile)
	dom.AddEventListener(`[data-action='logout']`, "click", comp.uiEventLogout)
	dom.AddEventListener(`[data-id='avatarCloseOverlay']`, "click", comp.uiEventToggleProfileClose)
}

func (comp *DashboardComponent) unbindUIEvents() {
	dom.RemoveEventListener(`[data-action='toggleProfile']`, "click", comp.uiEventLogout)
	dom.RemoveEventListener(`[data-action='logout']`, "click", comp.uiEventLogout)
	dom.RemoveEventListener(`[data-id='avatarCloseOverlay']`, "click", comp.uiEventLogout)
}

func (comp *DashboardComponent) eventMachine(machine *domain.Machine) {
	dom.ConsoleLog("Event machines", machine)

	//TODO: Parse data and show the template
	dashboardContainer := dom.QuerySelector(`[data-dock="dashboardContainer"]`)
	dashboardContainer.Set("innerHTML", "")

	machineComp := NewMachineComponent(comp.di)
	// TODO: Parse data and set fields here

	if err := machineComp.Render(); err != nil {
		dom.ConsoleLog("Machine comp render error", err.Error())
	}

	if err := machineComp.Mount(&dashboardContainer); err != nil {
		dom.ConsoleError("Machine comp mount error", err.Error())
		return
	}

}

func (comp *DashboardComponent) eventMachineError(err *orbital.ErrorResponse) {
	dom.ConsoleLog("Event machines error", err)
}

func (comp *DashboardComponent) uiEventToggleProfile(this js.Value, args []js.Value) interface{} {

	avatarMenu := dom.QuerySelector("[data-id='avatarMenu']")
	dom.RemoveClass(avatarMenu, "hidden")

	avatarCloseOverlay := dom.QuerySelector("[data-id='avatarCloseOverlay']")
	dom.RemoveClass(avatarCloseOverlay, "hidden")

	return nil
}

func (comp *DashboardComponent) uiEventToggleProfileClose(_ js.Value, _ []js.Value) interface{} {

	avatarMenu := dom.QuerySelector("[data-id='avatarMenu']")
	dom.AddClass(avatarMenu, "hidden")

	avatarCloseOverlay := dom.QuerySelector("[data-id='avatarCloseOverlay']")
	dom.AddClass(avatarCloseOverlay, "hidden")

	return nil
}

func (comp *DashboardComponent) uiEventLogout(_ js.Value, _ []js.Value) interface{} {

	userRepo := domain.NewRepository[*domain.User](comp.di.Storage(), domain.UserStorageKey)
	_ = userRepo.Remove()

	authRepo := domain.NewRepository[*domain.Auth](comp.di.Storage(), domain.AuthStorageKey)
	_ = authRepo.Remove()

	comp.di.State().Set("state:orbital:authenticated", false)
	comp.Unmount()
	return nil
}
