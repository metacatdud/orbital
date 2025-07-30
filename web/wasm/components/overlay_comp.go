package components

import (
	"fmt"
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/dom"
	"orbital/web/wasm/pkg/state"
	"strings"
	"syscall/js"
)

const (
	OverlayComponentRegKey RegKey = "overlayComponent"
)

type OverlayConfig struct {
	Child   Component
	Title   string
	Icon    string
	Actions []string
	Css     []string
	OnClose func()
}

type OverlayComponent struct {
	*BaseComponent
	child     Component
	container *js.Value
	state     *state.State

	title   string
	icon    string
	actions []string
	css     []string
	onClose func()
}

func NewOverlayComponent(di *orbital.Dependency, cfg OverlayConfig) *OverlayComponent {
	base := NewBaseComponent(di, OverlayComponentRegKey, "orbital/overlay/overlay")
	comp := &OverlayComponent{
		BaseComponent: base,
		state:         di.State,
		child:         cfg.Child,
		title:         cfg.Title,
		icon:          cfg.Icon,
		css:           cfg.Css,
		actions:       cfg.Actions,
		onClose:       cfg.OnClose,
	}

	comp.OnInit(comp.onInit)

	return comp
}

func (comp *OverlayComponent) ID() RegKey {
	if comp.child != nil {
		return OverlayComponentRegKey.WithExtra("-", comp.child.ID())
	}

	return OverlayComponentRegKey

}

func (comp *OverlayComponent) Mount(container *js.Value) error {

	if comp.onInit != nil {
		comp.onInit()
	}

	data := map[string]any{
		"id":        comp.ID(),
		"windowCss": strings.Join(comp.css, " "),
		"title":     comp.title,
		"icon":      comp.icon,
		"actions":   comp.actions,
	}

	dom.ConsoleLog("OverlayComponentDATA", data)

	html, err := comp.Render(data)
	if err != nil {
		return fmt.Errorf("OverlayComponent: render failed: %w", err)
	}

	el := dom.CreateElementFromString(html)
	dom.AppendChild(*container, el)
	comp.container = container

	// Thought the BaseComponent.element
	comp.element = &el
	comp.RegisterContainers()

	childContainer := comp.GetContainer("overlayBody")
	if !childContainer.IsNull() && comp.child != nil {
		if err = comp.child.Mount(&childContainer); err != nil {
			dom.ConsoleError("overlayBody container not found")
			return err
		}
	}

	for _, uiEvt := range comp.uiEventHandlers {
		dom.AddEventListener(uiEvt.selector, uiEvt.event, uiEvt.cb)
	}

	if comp.onMount != nil {
		comp.onMount()
	}

	return nil
}

func (comp *OverlayComponent) Unmount() error {
	if comp.child != nil {
		comp.child.Unmount()
	}

	dom.AddClass(*comp.container, "hide")
	return comp.BaseComponent.Unmount()
}

func (comp *OverlayComponent) Render(data map[string]any) (string, error) {
	return comp.BaseComponent.Render(data)
}

func (comp *OverlayComponent) onInit() {
	comp.bindUIEvents()
}

func (comp *OverlayComponent) bindUIEvents() {
	if hasAction(comp.actions, "close") {
		comp.AddEventHandler(`[data-action="closeOverlay"]`, "click", comp.uiEventOverlayClose)
	}

	if hasAction(comp.actions, "minimize") {
		// TODO: Add handler for minimize
	}

	if hasAction(comp.actions, "maximize") {
		// TODO: Add handler for maximize
	}
}

func (comp *OverlayComponent) uiEventOverlayClose(_ js.Value, args []js.Value) any {
	if comp.onClose != nil {
		comp.onClose()
	}
	_ = comp.Unmount()
	return nil
}

func hasAction(actions []string, action string) bool {
	for _, a := range actions {
		if a == action {
			return true
		}
	}
	return false
}
