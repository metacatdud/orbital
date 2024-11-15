package dom

import (
	"fmt"
	"strings"
	"syscall/js"
)

var templates []*Element

type Element struct {
	ID         string
	Obj        js.Value
	IsTemplate bool
}

// AddModuleTemplate loads a new template into the DOM with a unique ID
func AddModuleTemplate(id string, htmlBin []byte) error {

	// Construct the path and check if the template already exists
	if exists(id) {
		return fmt.Errorf("template already added: [%s]", id)
	}

	tmpNode := Document().Obj.Call("createElement", "template")
	tmpNode.Set("innerHTML", string(htmlBin))
	tmpNode.Set("dataset.id", id)

	obj := tmpNode.Get("content").
		Call("querySelectorAll", "[data-template]")

	for i := 0; i < obj.Length(); i++ {
		tpl := obj.Index(i)
		tID := tpl.Get("dataset").
			Get("template").
			String()

		el := &Element{
			ID:         tplId(id, tID),
			Obj:        tpl,
			IsTemplate: true,
		}

		templates = append(templates, el)
	}

	return nil
}

// GetTemplate retrieves and clones a specific template by ID
func GetTemplate(id string, templateId ...string) (js.Value, error) {

	tID := "default"
	if len(templateId) > 0 {
		tID = templateId[0]
	}

	// Find the template in the templates list by ID
	for _, tpl := range templates {
		if tpl.ID == tplId(id, tID) {
			// Query and clone the template content
			htmlNode := tpl.Obj.
				Get("firstElementChild").
				Call("cloneNode", true)

			return htmlNode, nil
		}
	}

	return js.Null(), fmt.Errorf("template does not exist: [%s]", id)
}

// exists checks if a template with a given ID already exists
func exists(id string) bool {
	for _, t := range templates {
		if t.ID == id {
			return true
		}
	}
	return false
}

func tplId(moduleID, tplID string) string {
	return strings.Join([]string{moduleID, tplID}, "/")
}
