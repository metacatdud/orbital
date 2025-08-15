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
	namespace := "dashboard/app/appLauncher"
	if len(data.Apps) > 0 {
		namespace = "dashboard/app/appsGroup"
	}

	base := NewBaseComponent(di, id, namespace)
	comp := &AppLauncherComponent{
		BaseComponent: base,
		id:            id,
		data:          data,
		events:        di.Events,
		state:         di.State,
	}

	comp.bindUIEvents()
	comp.OnMount(comp.onMountHandler)

	return comp
}

func (comp *AppLauncherComponent) ID() RegKey {
	return comp.id
}

func (comp *AppLauncherComponent) Mount(container *js.Value) error {
	if comp.onInit != nil {
		comp.onInit()
	}

	data := map[string]any{
		"id":    comp.id,
		"title": comp.Title(),
		"icon":  comp.Icon(),
		"data":  comp.data,
	}

	html, err := comp.Render(data)
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

func (comp *AppLauncherComponent) Render(data map[string]any) (string, error) {
	return comp.BaseComponent.Render(data)
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

	comp.AddEventHandler(
		fmt.Sprintf(`[data-id="%s"][data-action="launchGroup"]`, comp.ID()),
		"click",
		comp.uiEventLaunchGroup,
	)

	comp.AddEventHandler(
		fmt.Sprintf(`[data-target="%s"] [data-action="closeGroup"]`, comp.ID()),
		"click",
		comp.uiEventLaunchGroupClose,
	)
}

func (comp *AppLauncherComponent) onMountHandler() {
	if len(comp.data.Apps) > 0 {
		comp.setupChildApps()
	}
}

func (comp *AppLauncherComponent) uiEventLaunchApp(_ js.Value, args []js.Value) any {

	// TODO: Instantiate

	app := NewAppComponent(comp.DI, RegKey(comp.data.ID), comp.data)

	comp.events.Emit("evt:overlay:show", OverlayConfig{
		Child:   app,
		Title:   comp.Title(),
		Icon:    comp.Icon(),
		Actions: []string{"close", "minimize"},
		Css:     []string{"large"},
		OnClose: nil,
	})
	return nil
}

func (comp *AppLauncherComponent) uiEventLaunchGroup(_ js.Value, args []js.Value) any {
	e := args[0]
	e.Call("stopPropagation")

	targetSel := fmt.Sprintf(`[data-target="%s"]`, comp.ID())
	targetMenu := dom.QuerySelector(targetSel)

	dom.ToggleClass(targetMenu, "hide")
	return nil
}

func (comp *AppLauncherComponent) uiEventLaunchGroupClose(_ js.Value, args []js.Value) any {
	e := args[0]
	e.Call("stopPropagation")

	targetSel := fmt.Sprintf(`[data-target="%s"]`, comp.ID())
	targetMenu := dom.QuerySelector(targetSel)

	dom.AddClass(targetMenu, "hide")
	return nil
}

func (comp *AppLauncherComponent) setupChildApps() {

	container := comp.GetContainer("appsList")
	if container.IsNull() {
		dom.ConsoleError("appsList component container is null", comp.ID())
		return
	}

	dom.SetInnerHTML(container, "")

	for _, app := range comp.data.Apps {
		appLauncher := NewAppLauncherComponent(comp.DI, AppComponentRegKey.WithExtra("-", app.ID), app)
		appLauncher.Mount(&container)
	}
}
