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
}

func NewOverlayComponent(di *orbital.Dependency) *OverlayComponent {
	base := NewBaseComponent(di, OverlayComponentRegKey, "orbital/overlay/overlay")
	comp := &OverlayComponent{
		BaseComponent: base,
		state:         di.State,
	}

	return comp
}

func (comp *OverlayComponent) ID() RegKey {
	return OverlayComponentRegKey
}

func (comp *OverlayComponent) Mount(container *js.Value) error {

	// Check state and set default
	if rawState := comp.DI.State.Get("overlayChild"); rawState == nil {
		comp.state.Set("overlayChild", LoginComponentRegKey)
	}

	childName := comp.state.Get("overlayChild").(RegKey)

	// Lookup new component
	childComp, err := LookupComponent(childName, comp.DI)
	if err != nil {
		return err
	}

	// See if metadata is provided
	title, icon := "", ""
	if mp, ok := childComp.(MetaProvider); ok {
		title, icon = mp.Title(), mp.Icon()
	}

	// Render overlay with data
	html, err := comp.Render(map[string]interface{}{
		"title": title,
		"icon":  icon,
	})
	if err != nil {
		dom.ConsoleError("cannot render overlay", err.Error())
		return err
	}

	el := dom.CreateElementFromString(html)

	dom.AppendChild(*container, el)

	comp.element = &el
	comp.RegisterContainers()

	// Prepare and mount child
	comp.child = childComp
	overlayBody := comp.GetContainer("overlayBody")
	if overlayBody.IsNull() {
		return fmt.Errorf("dock area [overlayBody] not found")
	}

	comp.bindUIEvents()

	if err = childComp.Mount(&overlayBody); err != nil {
		return fmt.Errorf("cannot mount overlay component %s", childName)
	}

	comp.state.Set("state:overlay:toggle", true)

	// TODO: Redo this as it turns very ugly
	//comp.state.Watch("overlayChild", func(_, newValue any) {
	//	newV := newValue.(RegKey)
	//	dom.ConsoleLog("Change child", newV)
	//
	//	//// Unmount old child if any
	//	//if comp.child != nil {
	//	//	comp.child.Unmount()
	//	//}
	//	//
	//	//// Lookup new component
	//	//childComp, err = LookupComponent(newV, comp.DI)
	//	//if err != nil {
	//	//	dom.ConsoleError("cannot find overlay component", err.Error())
	//	//	return
	//	//}
	//	//
	//	//comp.child = childComp
	//	//
	//	//title, icon = "", ""
	//	//if mp, ok := childComp.(MetaProvider); ok {
	//	//	title, icon = mp.Title(), mp.Icon()
	//	//}
	//	//
	//	//// Render overlay with data
	//	//html, err = comp.Render(map[string]interface{}{
	//	//	"title": title,
	//	//	"icon":  icon,
	//	//})
	//	//if err != nil {
	//	//	dom.ConsoleError("cannot render overlay", err.Error())
	//	//	return
	//	//}
	//	//
	//	//el := dom.CreateElementFromString(html)
	//
	//})

	return nil
}

func (comp *OverlayComponent) Unmount() error {
	comp.state.Set("state:overlay:toggle", false)
	if comp.child != nil {
		comp.child.Unmount()
	}

	return comp.BaseComponent.Unmount()
}

func (comp *OverlayComponent) bindUIEvents() {
	dom.AddEventListener(`[data-action="closeOverlay"]`, "click", comp.uiEventOverlayClose)
}

func (comp *OverlayComponent) uiEventOverlayClose(_ js.Value, args []js.Value) interface{} {
	comp.Unmount()
	return nil
}

//TODO: Make components report this in
