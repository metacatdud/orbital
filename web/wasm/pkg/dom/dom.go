package dom

import (
	"strings"
	"syscall/js"
)

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

func CreateElementFromString(htmlStr string) js.Value {
	doc := js.Global().Get("document")
	tag := extractTagName(htmlStr)

	var container js.Value

	switch tag {
	case "tr", "td", "th":
		container = doc.Call("createElement", "tbody")
	case "option":
		container = doc.Call("createElement", "select")
	default:
		container = doc.Call("createElement", "div")
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
