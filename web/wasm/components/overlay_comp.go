package components

import (
	"fmt"
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/dom"
	"orbital/web/wasm/pkg/state"
	"syscall/js"
)

const (
	OverlayComponentRegKey RegKey = "overlayComponent"
)

type OverlayComponent struct {
	*BaseComponent
	child Component
	state *state.State

	titleText string
	iconClass string
	container js.Value
}

func NewOverlayComponent(di *orbital.Dependency) *OverlayComponent {
	base := NewBaseComponent(di, OverlayComponentRegKey, "orbital/overlay/overlay")
	comp := &OverlayComponent{
		BaseComponent: base,
		state:         di.State,
	}

	comp.OnMount(comp.onMountHandler)

	di.State.Watch("state:overlay:currentChild", func(oldVal, newVal interface{}) {
		if oldVal == nil {
			return
		}

		if oldVal == newVal {
			return
		}

		newKey, ok := newVal.(RegKey)
		if !ok {
			dom.ConsoleError("Overlay: state:overlay:currentChild changed to non-RegKey")
			return
		}

		comp.swapChild(newKey)
	})

	comp.bindUIEvents()

	return comp
}

func (comp *OverlayComponent) ID() RegKey {
	return OverlayComponentRegKey
}

func (comp *OverlayComponent) Mount(container *js.Value) error {

	if comp.onInit != nil {
		comp.onInit()
	}

	raw := comp.state.Get("state:overlay:currentChild")
	if raw == nil {
		dom.ConsoleError("no child component set for overlay")
		return comp.Unmount()
	}

	key, _ := raw.(RegKey)
	dom.ConsoleLog("overlay mount with child", key)

	child, err := LookupComponent(key, comp.DI)
	if err != nil {
		return fmt.Errorf("overlay: failed to lookup child %q: %w", key, err)
	}

	comp.child = child

	comp.titleText = ""
	comp.iconClass = ""
	if mp, ok := child.(MetaProvider); ok {
		comp.titleText = mp.Title()
		comp.iconClass = mp.Icon()
	}

	html, err := comp.Render(nil)
	if err != nil {
		return err
	}

	el := dom.CreateElementFromString(html)
	dom.AppendChild(*container, el)

	comp.element = &el
	comp.RegisterContainers()

	comp.container = *container

	if comp.onMount != nil {
		comp.onMount()
	}
	return nil
}

func (comp *OverlayComponent) Unmount() error {
	comp.state.Set("state:overlay:toggle", false)
	if comp.child != nil {
		comp.child.Unmount()
	}

	return comp.BaseComponent.Unmount()
}

func (comp *OverlayComponent) Render(_ map[string]interface{}) (string, error) {
	return comp.BaseComponent.Render(map[string]interface{}{
		"title": comp.titleText,
		"icon":  comp.iconClass,
	})
}

func (comp *OverlayComponent) onMountHandler() {

	dom.ConsoleLog("overlay mount child", comp.child.ID())

	container := comp.GetContainer("overlayBody")
	if container.IsNull() {
		dom.ConsoleError("overlayBody container not found")
		return
	}
	if comp.child != nil {
		comp.child.Mount(&container)
	}

	comp.state.Set("state:overlay:toggle", true)
}

func (comp *OverlayComponent) bindUIEvents() {
	comp.AddEventHandler(`[data-action="closeOverlay"]`, "click", comp.uiEventOverlayClose)
}

func (comp *OverlayComponent) uiEventOverlayClose(_ js.Value, args []js.Value) interface{} {
	comp.Unmount()
	return nil
}

func (comp *OverlayComponent) swapChild(key RegKey) {

	container := comp.container
	if container.IsNull() {
		dom.ConsoleError("cannot swapChild: container is null")
		return
	}

	container.Set("innerHTML", "")
	dom.RemoveElement(*comp.element)

	if comp.child != nil {
		comp.child.Unmount()
	}

	newChild, err := LookupComponent(key, comp.DI)
	if err != nil {
		dom.ConsoleError("overlay:swapChild", err.Error())
		return
	}
	comp.child = newChild

	comp.titleText = ""
	comp.iconClass = ""
	if mp, ok := newChild.(MetaProvider); ok {
		comp.titleText = mp.Title()
		comp.iconClass = mp.Icon()
	}

	html, err := comp.Render(nil)
	if err != nil {
		dom.ConsoleError("overlay update failed", err.Error())
		return
	}

	el := dom.CreateElementFromString(html)

	comp.element = &el
	dom.AppendChild(container, el)

	comp.RegisterContainers()
	comp.bindUIEvents()

	for _, ob := range comp.observers {
		ob.OnParentRender()
	}

	body := comp.GetContainer("overlayBody")
	if container.IsNull() {
		dom.ConsoleError("overlayBody container not found")
		return
	}
	newChild.Mount(&body)
	comp.state.Set("state:overlay:toggle", true)
}
