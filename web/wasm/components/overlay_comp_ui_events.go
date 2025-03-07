package components

import (
	"orbital/web/wasm/pkg/dom"
	"syscall/js"
)

func (comp *OverlayComponent) bindUIEvents() {
	dom.AddEventListener(`[data-action="startOrbital"]`, "click", comp.uiEventHandleClose)
}

func (comp *OverlayComponent) unbindUIEvents() {}

func (comp *OverlayComponent) uiEventHandleClose(_ js.Value, _ []js.Value) interface{} {
	currentComp := comp.state.Get("state:overlay:activeComponent").(string)
	if currentComp == "login" {
		dom.ConsoleLog("Well if we close that than what?")
		return nil
	}

	return nil
}
