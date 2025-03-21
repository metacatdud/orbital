package components

import (
	"orbital/web/wasm/pkg/dom"
	"syscall/js"
)

func (comp *TaskbarComponent) bindUIEvents() {
	dom.AddEventListener(`[data-action="startOrbital"]`, "click", comp.uiEventStartOrbital)
	dom.Document().Call("addEventListener", "click", js.FuncOf(comp.uiEventStartOrbitalHide))

}

func (comp *TaskbarComponent) unbindUIEvents() {
	dom.RemoveEventListener(`[data-action="startOrbital"]`, "click", comp.uiEventStartOrbital)
	dom.Document().Call("removeEventListener", "click", js.FuncOf(comp.uiEventStartOrbitalHide))
}

func (comp *TaskbarComponent) uiEventStartOrbital(_ js.Value, args []js.Value) interface{} {
	e := args[0]
	e.Call("stopPropagation")

	startMenu := dom.QuerySelector(`[data-id="startMenu"]`)
	dom.ToggleClass(startMenu, "hide")

	return nil
}

// uiEventStartOrbitalHide close menu if clicking outside
func (comp *TaskbarComponent) uiEventStartOrbitalHide(this js.Value, args []js.Value) interface{} {
	e := args[0]
	target := e.Get("target")

	startMenu := dom.QuerySelector(`[data-id="startMenu"]`)
	startButton := dom.QuerySelector(`[data-action="startOrbital"]`)

	if !startMenu.Call("contains", target).Bool() && !target.Equal(startButton) {
		dom.AddClass(startMenu, "hide")
	}

	return nil
}
