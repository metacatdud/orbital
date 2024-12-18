package dom

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

var (
	document *Element
	dom      *Element
)

func init() {
	dom = &Element{
		ID:         "dom",
		Obj:        js.Global(),
		IsTemplate: false,
	}

	document = &Element{
		ID:         "document",
		Obj:        dom.Obj.Get("document"),
		IsTemplate: false,
	}
}

// Callback represents a function for handling JavaScript events in Go
//type Callback func(js.Value, []js.Value) interface{}

// Document returns the root document element
func Document() *Element {
	return document
}

// DOM returns the root DOM object
//func DOM() *Element {
//	return dom
//}

// DocQuerySelector selects an element from the document by CSS selector
func DocQuerySelector(selector string) js.Value {
	return document.Obj.Call("querySelector", selector)
}

// DocQuerySelectorValue selects an attribute from an element in the document
func DocQuerySelectorValue(selector, val string) js.Value {
	return DocQuerySelector(selector).Get(val)
}

// ElementQuerySelect selects a child element from a given element by CSS selector
func ElementQuerySelect(el js.Value, selector string) js.Value {
	return el.Call("querySelector", selector)
}

// ElementQuerySelectValue retrieves a property from a selected child element
func ElementQuerySelectValue(el js.Value, selector, val string) js.Value {
	return ElementQuerySelect(el, selector).Get(val)
}

// Show makes an element visible by setting display to "block"
func Show(selector string) {
	DocQuerySelectorValue(selector, "style").Set("display", "block")
}

// Hide hides an element by setting display to "none"
func Hide(selector string) {
	DocQuerySelectorValue(selector, "style").Set("display", "none")
}

// SetInnerHTML sets the inner HTML content of an element selected by CSS selector
//func SetInnerHTML(selector, html string) {
//	DocQuerySelector(selector).Set("innerHTML", html)
//}

// AppendChild appends a new child element to a target element specified by selector
func AppendChild(parentSelector string, child js.Value) {
	DocQuerySelector(parentSelector).Call("appendChild", child)
}

// Clear removes all child elements of an element specified by selector
func Clear(selector string) {
	parent := DocQuerySelector(selector)
	for parent.Get("firstChild").Truthy() {
		parent.Call("removeChild", parent.Get("firstChild"))
	}
}

func PrintToConsole(data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	// Parse JSON data into JavaScript object
	jsObject := js.Global().Get("JSON").Call("parse", string(jsonData))
	js.Global().Get("console").Call("log", jsObject)
}
