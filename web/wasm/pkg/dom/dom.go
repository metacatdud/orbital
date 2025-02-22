package dom

import "syscall/js"

func QuerySelector(selector string) js.Value {
	elem := js.Global().Get("document").Call("querySelector", selector)
	if elem.Truthy() {
		return elem
	}
	return js.Null()
}

func QuerySelectorAll(selector string) []js.Value {
	nodeList := js.Global().Get("document").Call("querySelectorAll", selector)
	length := nodeList.Length()

	elements := make([]js.Value, length)
	for i := 0; i < length; i++ {
		elements[i] = nodeList.Index(i)
	}

	return elements
}

func CreateElement(tag string) js.Value {
	return js.Global().Get("document").Call("createElement", tag)
}

func SetInnerHTML(element js.Value, html string) {
	if element.Truthy() {
		element.Set("innerHTML", html)
	}
}

func AppendChild(parent js.Value, child js.Value) {
	if parent.Truthy() && child.Truthy() {
		parent.Call("appendChild", child)
	}
}

func RemoveElement(element js.Value) {
	if element.Truthy() {
		element.Call("remove")
	}
}

func GetValue(selector string) string {
	elem := QuerySelector(selector)
	if elem.IsNull() {
		return ""
	}
	return elem.Get("value").String()
}

func GetPlaceholder(selector string) js.Value {
	elem := QuerySelector(selector)
	if elem.IsNull() {
		return js.Null()
	}
	return elem
}

func AddEventListener(selector, event string, callback func(js.Value, []js.Value) interface{}) {
	elem := QuerySelector(selector)
	if elem.IsNull() {
		return
	}
	handler := js.FuncOf(callback)
	elem.Call("addEventListener", event, handler)
}

func RemoveEventListener(selector, event string, callback func(js.Value, []js.Value) interface{}) {
	elem := QuerySelector(selector)
	if elem.IsNull() {
		return
	}
	handler := js.FuncOf(callback)
	elem.Call("removeEventListener", event, handler)
}

func AddClass(element js.Value, className string) {
	if element.Truthy() {
		element.Get("classList").Call("add", className)
	}
}

func RemoveClass(element js.Value, className string) {
	if element.Truthy() {
		element.Get("classList").Call("remove", className)
	}
}
