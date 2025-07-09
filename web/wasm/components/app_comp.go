package components

import (
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/dom"
	"orbital/web/wasm/pkg/state"
	"orbital/web/wasm/service"
	"syscall/js"
)

const (
	AppComponentRegKey RegKey = "application"
)

type AppComponent struct {
	*BaseComponent
	id    RegKey
	data  service.App
	state *state.State
}

var _ MetaProvider = (*AppComponent)(nil)

func NewAppComponent(di *orbital.Dependency, id RegKey, data service.App) *AppComponent {
	base := NewBaseComponent(di, id, "dashboard/app/appLauncher")
	comp := &AppComponent{
		BaseComponent: base,
		id:            id,
		data:          data,
		state:         di.State,
	}

	comp.bindUIEvents()

	return comp
}

func (comp *AppComponent) ID() RegKey {
	return comp.id
}

func (comp *AppComponent) Mount(container *js.Value) error {
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

func (comp *AppComponent) Unmount() error {
	return comp.BaseComponent.Unmount()
}

func (comp *AppComponent) Render(_ map[string]interface{}) (string, error) {
	return comp.BaseComponent.Render(map[string]interface{}{
		"title": comp.Title(),
		"icon":  comp.Icon(),
	})
}

func (comp *AppComponent) Title() string {
	return comp.data.Name
}

func (comp *AppComponent) Icon() string {
	return comp.data.Icon
}

func (comp *AppComponent) bindUIEvents() {
	comp.AddEventHandler(`[data-action="launchApp"]`, "click", comp.uiEventLaunchApp)
}

func (comp *AppComponent) uiEventLaunchApp(_ js.Value, args []js.Value) interface{} {
	dom.ConsoleLog("Open overlay with", comp.ID())
	comp.state.Set("state:overlay:toggle", true)
	return nil
}
