package components

import (
	"fmt"
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/dom"
	"orbital/web/wasm/pkg/events"
	"orbital/web/wasm/pkg/state"
	"orbital/web/wasm/service"
	"syscall/js"
)

const (
	AppComponentRegKey RegKey = "application"
)

type AppLauncherComponent struct {
	*BaseComponent
	id     RegKey
	data   service.App
	events *events.Event
	state  *state.State
}

var _ MetaProvider = (*AppLauncherComponent)(nil)

func NewAppLauncherComponent(di *orbital.Dependency, id RegKey, data service.App) *AppLauncherComponent {
	base := NewBaseComponent(di, id, "dashboard/app/appLauncher")
	comp := &AppLauncherComponent{
		BaseComponent: base,
		id:            id,
		data:          data,
		events:        di.Events,
		state:         di.State,
	}

	comp.bindUIEvents()

	return comp
}

func (comp *AppLauncherComponent) ID() RegKey {
	return comp.id
}

func (comp *AppLauncherComponent) Mount(container *js.Value) error {
	if comp.onInit != nil {
		comp.onInit()
	}

	html, err := comp.Render(nil)
	if err != nil {
		return err
	}

	el := dom.CreateElementFromString(html)
	dom.AppendChild(*container, el)

	comp.element = &el
	comp.RegisterContainers()

	for _, uiEvt := range comp.uiEventHandlers {
		dom.AddEventListener(uiEvt.selector, uiEvt.event, uiEvt.cb)
	}

	if comp.onMount != nil {
		comp.onMount()
	}

	return nil
}

func (comp *AppLauncherComponent) Unmount() error {
	return comp.BaseComponent.Unmount()
}

func (comp *AppLauncherComponent) Render(_ map[string]any) (string, error) {
	return comp.BaseComponent.Render(map[string]any{
		"id":    comp.ID().String(),
		"title": comp.Title(),
		"icon":  comp.Icon(),
	})
}

func (comp *AppLauncherComponent) Title() string {
	return comp.data.Name
}

func (comp *AppLauncherComponent) Icon() string {
	return comp.data.Icon
}

func (comp *AppLauncherComponent) bindUIEvents() {
	comp.AddEventHandler(
		fmt.Sprintf(`[data-id="%s"][data-action="launchApp"]`, comp.ID()),
		"click",
		comp.uiEventLaunchApp,
	)
}

func (comp *AppLauncherComponent) uiEventLaunchApp(_ js.Value, args []js.Value) any {

	// TODO: Instantiate

	app := NewAppComponent(comp.DI, RegKey(comp.data.ID), comp.data)

	comp.events.Emit("evt:overlay:show", OverlayConfig{
		Child:   app,
		Title:   comp.Title(),
		Icon:    comp.Icon(),
		Actions: []string{"close"},
		Css:     []string{"large"},
		OnClose: nil,
	})
	return nil
}
