package dom

import (
	"strings"
	"syscall/js"
)

func Document() js.Value {
	return js.Global().Get("document")
}

func QuerySelector(selector string) js.Value {
	elem := Document().Call("querySelector", selector)
	if elem.Truthy() {
		return elem
	}
	return js.Null()
}

func QuerySelectorAll(selector string) []js.Value {
	nodeList := Document().Call("querySelectorAll", selector)
	length := nodeList.Length()

	elements := make([]js.Value, length)
	for i := 0; i < length; i++ {
		elements[i] = nodeList.Index(i)
	}

	return elements
}

func QuerySelectorFromElement(parent js.Value, selector string) js.Value {
	if parent.IsNull() {
		return js.Null()
	}

	elem := parent.Call("querySelector", selector)
	if elem.Truthy() {
		return elem
	}
	return js.Null()
}

func QuerySelectorAllFromElement(parent js.Value, selector string) []js.Value {
	if parent.IsNull() {
		return []js.Value{}
	}

	nodeList := parent.Call("querySelectorAll", selector)
	length := nodeList.Length()

	elements := make([]js.Value, length)
	for i := 0; i < length; i++ {
		elements[i] = nodeList.Index(i)
	}

	return elements
}

func CreateElement(tag string) js.Value {
	return Document().Call("createElement", tag)
}

func CreateElementFromString(htmlStr string) js.Value {
	tag := extractTagName(htmlStr)

	var container js.Value

	switch tag {
	case "tr", "td", "th":
		container = Document().Call("createElement", "tbody")
	case "option":
		container = Document().Call("createElement", "select")
	default:
		container = Document().Call("createElement", "div")
	}

	container.Set("innerHTML", htmlStr)
	return container.Get("firstElementChild")
}

func extractTagName(s string) string {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "<") {
		return ""
	}
	s = s[1:]
	end := strings.IndexAny(s, " >")
	if end == -1 {
		return ""
	}
	return s[:end]
}

func SetInnerHTML(element js.Value, html string) {
	if element.Truthy() {
		element.Set("innerHTML", html)
	}
}

func AppendChild(parent, child js.Value) {
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

func ToggleClass(element js.Value, className string) {
	if element.Truthy() {
		element.Get("classList").Call("toggle", className)
	}
}
