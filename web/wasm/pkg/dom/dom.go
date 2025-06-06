package dom

import (
	"fmt"
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

func AddEventListener(selector, event string, callback js.Func) {
	elem := QuerySelector(selector)
	if elem.IsNull() {
		return
	}
	elem.Call("addEventListener", event, callback)
}

func RemoveEventListener(selector, event string, callback js.Func) {
	elem := QuerySelector(selector)
	if elem.IsNull() {
		return
	}
	elem.Call("removeEventListener", event, callback)
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

func GetValue(attr, key string) string {
	sel := fmt.Sprintf("[data-%s='%s']", attr, key)
	el := QuerySelector(sel)
	if el.IsNull() || el.IsUndefined() {
		return ""
	}

	// Try getting value from input, select, textarea
	if v := el.Get("value"); v.Type() == js.TypeString {
		// Handle radio buttons
		if el.Get("type").String() == "radio" {
			if !el.Get("checked").Bool() {
				return ""
			}
		}
		return v.String()
	}

	// Fallback to textContent
	if t := el.Get("textContent"); t.Type() == js.TypeString {
		return strings.TrimSpace(t.String())
	}

	return ""
}

func SetValue(attr, key, val string) {
	sel := fmt.Sprintf("[data-%s='%s']", attr, key)
	el := QuerySelector(sel)
	if el.IsNull() || el.IsUndefined() {
		ConsoleError("selector not found", sel)
		return
	}

	switch el.Get("nodeName").String() {
	case "INPUT":
		typ := el.Get("type").String()
		switch typ {
		case "checkbox":
			realVal := strings.ToLower(val)
			el.Set("checked", realVal == "true" || realVal == "1" || realVal == "on")
		case "radio":
			el.Set("checked", true)
		default:
			el.Set("value", val)
		}
	case "SELECT":
		el.Set("value", val)
	case "OPTION":
		el.Set("value", val)

		if el.Get("value").String() == val {
			el.Set("selected", true)
		}
	default:
		el.Set("textContent", val)
	}
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
