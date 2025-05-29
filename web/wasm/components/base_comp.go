package components

import (
	"bytes"
	"html/template"
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/dom"
	"syscall/js"
)

type BaseComponent struct {
	DI        *orbital.Dependency
	id        RegKey
	namespace string
	element   *js.Value
	docks     map[string]js.Value
	tpl       *template.Template

	onInit, onMount, onUpdate, onUnmount func()

	uiEventHandlers []struct {
		selector string
		event    string
		cb       func(js.Value, []js.Value) interface{}
	}

	unwatchFns []func()
	observers  []ParentRenderObserver
}

func NewBaseComponent(di *orbital.Dependency, id RegKey, ns string) *BaseComponent {
	return &BaseComponent{
		DI:        di,
		id:        id,
		namespace: ns,
		docks:     make(map[string]js.Value),
	}
}

func (comp *BaseComponent) ID() RegKey {
	return comp.id
}

func (comp *BaseComponent) OnInit(fn func()) {
	comp.onInit = fn
}

func (comp *BaseComponent) OnMount(fn func()) {
	comp.onMount = fn
}

func (comp *BaseComponent) OnUpdate(fn func()) {
	comp.onUpdate = fn
}

func (comp *BaseComponent) OnUnmount(fn func()) {
	comp.onUnmount = fn
}

func (comp *BaseComponent) AddEventHandler(sel, evt string, cb func(js.Value, []js.Value) interface{}) {
	comp.uiEventHandlers = append(comp.uiEventHandlers,
		struct {
			selector string
			event    string
			cb       func(js.Value, []js.Value) interface{}
		}{
			selector: sel,
			event:    evt,
			cb:       cb,
		})
}

func (comp *BaseComponent) Watch(key string, cb func(oldV, newV interface{})) {
	unwatchFn := comp.DI.State.Watch(key, func(oldValue, newValue interface{}) {
		comp.update()

		if comp.onUpdate != nil {
			comp.onUpdate()
		}

		cb(oldValue, newValue)
	})
	comp.unwatchFns = append(comp.unwatchFns, unwatchFn)
}

func (comp *BaseComponent) RegisterObserver(o ParentRenderObserver) {
	comp.observers = append(comp.observers, o)
}

func (comp *BaseComponent) RegisterContainers() {
	if comp.element == nil {
		return
	}

	dockingAreas := dom.QuerySelectorAllFromElement(*comp.element, `[data-dock]`)
	for _, d := range dockingAreas {
		dockName := d.Get("dataset").Get("dock").String()
		comp.docks[dockName] = d
	}
}

func (comp *BaseComponent) GetContainer(name string) js.Value {
	if dockEl, ok := comp.docks[name]; ok && dockEl.Truthy() {
		return dockEl
	}

	return js.Null()
}

func (comp *BaseComponent) Render(data map[string]interface{}) (string, error) {
	if comp.tpl == nil {
		tpl, err := comp.DI.Templates.Get(comp.namespace)
		if err != nil {
			return "", err
		}

		comp.tpl = tpl
	}

	var buf bytes.Buffer
	if data == nil {
		data = comp.DI.State.GetAll()
	}

	if err := comp.tpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (comp *BaseComponent) Mount(container *js.Value) error {
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

func (comp *BaseComponent) Unmount() error {

	// Cleanup state
	for _, unwatchFn := range comp.unwatchFns {
		unwatchFn()
	}

	comp.unwatchFns = nil

	// Unmount component
	if comp.onUnmount != nil {
		comp.onUnmount()
	}

	// Cleanup DOM
	if comp.element != nil {
		dom.RemoveElement(*comp.element)
		comp.element = nil
	}

	comp.docks = nil

	// Cleanup DOM events
	for _, uiEvt := range comp.uiEventHandlers {
		dom.RemoveEventListener(uiEvt.selector, uiEvt.event, uiEvt.cb)
	}

	comp.uiEventHandlers = nil

	return nil
}

func (comp *BaseComponent) update() {
	if comp.element == nil {
		return
	}

	html, err := comp.Render(nil)
	if err != nil {
		dom.ConsoleError("cannot update", err.Error())
		return
	}

	(*comp.element).Set("innerHTML", html)
	for _, ob := range comp.observers {
		ob.OnParentRender()
	}
}
