package components

import (
	"orbital/web/wasm/pkg/dom"
	"syscall/js"
)

func (comp *TaskbarStartComponent) bindUIEvents() {
	dom.AddEventListener(`[data-action="startOrbital"]`, "click", comp.uiEventStartOrbital)
	dom.Document().Call("addEventListener", "click", js.FuncOf(comp.uiEventStartOrbitalHide))

}

func (comp *TaskbarStartComponent) unbindUIEvents() {
	dom.RemoveEventListener(`[data-action="startOrbital"]`, "click", comp.uiEventStartOrbital)
	dom.Document().Call("removeEventListener", "click", js.FuncOf(comp.uiEventStartOrbitalHide))
}

func (comp *TaskbarStartComponent) uiEventStartOrbital(_ js.Value, args []js.Value) interface{} {
	e := args[0]
	e.Call("stopPropagation")

	startMenu := dom.QuerySelector(`[data-id="startMenu"]`)
	dom.ToggleClass(startMenu, "hide")
	return nil
}

// uiEventStartOrbitalHide close menu if clicking outside
func (comp *TaskbarStartComponent) uiEventStartOrbitalHide(this js.Value, args []js.Value) interface{} {
	e := args[0]
	target := e.Get("target")

	dom.ConsoleLog("target", target)

	startMenu := dom.QuerySelector(`[data-id="startMenu"]`)
	if !startMenu.Call("contains", target).Bool() && !target.Call("contains", startMenu).Bool() {
		dom.AddClass(startMenu, "hide")
	}
	return nil
}
