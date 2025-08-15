package components

import (
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/dom"
	"syscall/js"
)

const (
	TaskbarComponentRegKey RegKey = "taskbarComponent"
)

type TaskbarComponent struct {
	*BaseComponent
	taskbarStartComp Component
}

func NewTaskbarComponent(di *orbital.Dependency) *TaskbarComponent {

	base := NewBaseComponent(di, TaskbarComponentRegKey, "orbital/taskbar/taskbar")

	comp := &TaskbarComponent{BaseComponent: base}
	comp.bindUIEvents()

	base.OnMount(comp.onMountHandler)
	base.OnUnmount(comp.onUnmountHandler)

	return comp
}

func (comp *TaskbarComponent) ID() RegKey {
	return TaskbarComponentRegKey
}

func (comp *TaskbarComponent) onMountHandler() {
	comp.mountTaskbarStartComp()

	// TODO: Find a better place to add events which are outside component's context
	dom.Document().Call("addEventListener", "click", js.FuncOf(comp.uiEventStartOrbitalHide))
}

func (comp *TaskbarComponent) onUnmountHandler() {
	dom.Document().Call("removeEventListener", "click", js.FuncOf(comp.uiEventStartOrbitalHide))
}

func (comp *TaskbarComponent) mountTaskbarStartComp() {
	container := comp.GetContainer("startMenu")
	if container.IsNull() {
		dom.ConsoleError("taskbar component container is null", comp.ID())
		return
	}

	dom.SetInnerHTML(container, "")

	taskbarStartComp := NewTaskbarStartComp(comp.DI)
	comp.taskbarStartComp = taskbarStartComp

	if err := taskbarStartComp.Mount(&container); err != nil {
		dom.ConsoleError("taskbar component cannot be mounted", err.Error(), comp.ID())
		return
	}
}

func (comp *TaskbarComponent) bindUIEvents() {
	comp.AddEventHandler(`[data-action="startOrbital"]`, "click", comp.uiEventStartOrbital)
}

func (comp *TaskbarComponent) uiEventStartOrbital(_ js.Value, args []js.Value) any {
	e := args[0]
	e.Call("stopPropagation")

	startMenu := dom.QuerySelector(`[data-id="startMenu"]`)
	dom.ToggleClass(startMenu, "hide")

	return nil
}

// uiEventStartOrbitalHide close menu if clicking outside
func (comp *TaskbarComponent) uiEventStartOrbitalHide(_ js.Value, args []js.Value) interface{} {
	e := args[0]
	target := e.Get("target")

	startMenu := dom.QuerySelector(`[data-id="startMenu"]`)
	startButton := dom.QuerySelector(`[data-action="startOrbital"]`)

	if !startMenu.Call("contains", target).Bool() && !target.Equal(startButton) {
		dom.AddClass(startMenu, "hide")
	}

	return nil
}
