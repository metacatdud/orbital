package components

import (
	"bytes"
	"errors"
	"orbital/web/w2/orbital"
	"orbital/web/w2/pkg/dom"
	"syscall/js"
)

const (
	TaskbarCompKey = "taskbarComponent"
)

type TaskbarComp struct {
	BaseComp
	di      *orbital.Dependency
	docks   map[string]js.Value
	element js.Value

	taskbarStartComp   orbital.Component
	taskbarAppsComp    orbital.Component
	taskbarSysTrayComp orbital.Component
}

var _ orbital.ContainerComponent = (*TaskbarComp)(nil)

func NewTaskbarComp(di *orbital.Dependency) *TaskbarComp {
	comp := &TaskbarComp{
		BaseComp: BaseComp{docks: make(map[string]js.Value)},
		di:       di,
	}

	return comp
}

func (comp *TaskbarComp) ID() string {
	return TaskbarCompKey
}

func (comp *TaskbarComp) Namespace() string {
	return "orbital/taskbar/taskbar"
}

func (comp *TaskbarComp) Mount(container *js.Value) error {
	if !container.Truthy() {
		return errors.New("container does not exist")
	}

	if comp.element.IsNull() {
		return errors.New("element is missing")
	}

	dom.AppendChild(*container, comp.element)

	comp.mountTaskbarStart()
	comp.bindUIEvents()

	return nil
}

func (comp *TaskbarComp) Unmount() error {
	if !comp.element.IsNull() {
		dom.RemoveElement(comp.element)
		comp.element = js.Null()
	}

	comp.unbindUIEvents()

	return nil
}

func (comp *TaskbarComp) Render() error {
	tpl, err := comp.di.TplRegistry().Get(comp.Namespace())
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, nil); err != nil {
		return err
	}

	comp.element = dom.CreateElementFromString(buf.String())

	comp.SetContainers(comp.element)

	return nil
}

func (comp *TaskbarComp) mountTaskbarStart() {
	container := comp.GetContainer("startMenu")
	if container.IsNull() {
		dom.ConsoleError("taskbarStartComp component container is null", comp.ID())
		return
	}
	dom.SetInnerHTML(container, "")

	var err error
	comp.taskbarStartComp, err = orbital.Lookup[*TaskbarStartComp](TaskbarStartCompKey)
	if err != nil {
		dom.ConsoleError("taskbarStartComp component cannot be created", comp.ID())
		return
	}

	if err = comp.taskbarStartComp.Render(); err != nil {
		dom.ConsoleError("taskbarStartComp component cannot be rendered", err.Error(), comp.ID())
		return
	}

	if err = comp.taskbarStartComp.Mount(&container); err != nil {
		dom.ConsoleError("taskbarStartComp component cannot be mounted", err.Error(), comp.ID())
		return
	}
}

func (comp *TaskbarComp) bindUIEvents() {
	dom.AddEventListener(`[data-action="startOrbital"]`, "click", comp.uiEventStartOrbital)
	dom.Document().Call("addEventListener", "click", js.FuncOf(comp.uiEventStartOrbitalHide))
}

func (comp *TaskbarComp) unbindUIEvents() {
	dom.RemoveEventListener(`[data-action="startOrbital"]`, "click", comp.uiEventStartOrbital)
	dom.Document().Call("removeEventListener", "click", js.FuncOf(comp.uiEventStartOrbitalHide))
}

func (comp *TaskbarComp) uiEventStartOrbital(_ js.Value, args []js.Value) interface{} {
	e := args[0]
	e.Call("stopPropagation")

	startMenu := dom.QuerySelector(`[data-id="startMenu"]`)
	dom.ToggleClass(startMenu, "hide")

	return nil
}

// uiEventStartOrbitalHide close menu if clicking outside
func (comp *TaskbarComp) uiEventStartOrbitalHide(this js.Value, args []js.Value) interface{} {
	e := args[0]
	target := e.Get("target")

	startMenu := dom.QuerySelector(`[data-id="startMenu"]`)
	startButton := dom.QuerySelector(`[data-action="startOrbital"]`)

	if !startMenu.Call("contains", target).Bool() && !target.Equal(startButton) {
		dom.AddClass(startMenu, "hide")
	}

	return nil
}
