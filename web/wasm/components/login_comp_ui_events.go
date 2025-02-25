package components

import (
	"orbital/web/wasm/pkg/dom"
	"syscall/js"
)

func (comp *LoginComponent) bindUIEvents() {
	dom.AddEventListener(`[data-action="login"]`, "click", comp.uiEventHandleLogin)
}

func (comp *LoginComponent) unbindUIEvents() {
	dom.RemoveEventListener(`[data-action="login"]`, "click", comp.uiEventHandleLogin)
}

func (comp *LoginComponent) uiEventHandleLogin(_ js.Value, _ []js.Value) interface{} {
	skInput := dom.GetValue("#privateKey")

	comp.di.Events().Emit("evt:auth:login:request", skInput)

	return nil
}
