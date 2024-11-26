package dom

import (
	"errors"
	"fmt"
	"strings"
	"syscall/js"
)

var elements []*Element

type Element struct {
	ID         string
	Obj        js.Value
	IsTemplate bool
}

func (el *Element) AppendChild(child *Element) error {
	if !el.Obj.Truthy() {
		return fmt.Errorf("element is not initialized")
	}

	if el.IsTemplate {
		return errors.New("cannot append child to a template element")
	}

	el.Obj.Call("appendChild", child.Obj)
	return nil
}

func (el *Element) CloneFromTemplate() *Element {
	htmlEl := el.Obj.Get("firstElementChild").Call("cloneNode", true)
	return &Element{
		ID:         el.ID,
		IsTemplate: false,
		Obj:        htmlEl,
	}
}

func (el *Element) Clone() *Element {
	htmlEl := el.Obj.Call("cloneNode", true)
	return &Element{
		ID:         el.ID,
		IsTemplate: false,
		Obj:        htmlEl,
	}
}

func (el *Element) Find(selector string) (*Element, error) {
	child := el.Obj.Call("querySelector", selector)
	if !child.Truthy() {
		return nil, errors.New("failed to find element")
	}

	childID := child.Get("dataset").Get("id").String()
	for _, elem := range elements {
		if elem.Obj.Equal(child) {
			return elem, nil
		}
	}

	// Create if not exist
	elem := &Element{
		ID:         tplId(el.ID, childID),
		Obj:        child,
		IsTemplate: false,
	}

	elements = append(elements, elem)

	return elem, nil
}

func (el *Element) Clean() (*Element, error) {
	if !el.Obj.Truthy() {
		return nil, fmt.Errorf("element is not initialized")
	}

	copyEl := el.Obj
	// Remove all child nodes
	for copyEl.Get("firstChild").Truthy() {
		copyEl.Call("removeChild", copyEl.Get("firstChild"))
	}

	return &Element{
		ID:         el.ID,
		Obj:        copyEl,
		IsTemplate: false,
	}, nil
}

// RegisterElement loads a new template into the DOM with a unique ID
func RegisterElement(id string, htmlBin []byte) error {

	// Construct the path and check if the template already exists
	if exists(id) {
		return fmt.Errorf("element already added: [%s]", id)
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

		fmt.Printf("Registering: %s\n", tplId(id, tID))

		elements = append(elements, el)
	}

	return nil
}

// GetElement retrieves and clones a specific template by ID
func GetElement(id string, templateId ...string) (*Element, error) {

	tID := "default"
	if len(templateId) > 0 {
		tID = templateId[0]
	}

	fmt.Println("GET:", tplId(id, tID))

	// Find the template in the elements list by ID
	for _, tpl := range elements {
		if tpl.ID == tplId(id, tID) {
			// Query and clone the template content
			return tpl, nil
		}
	}

	return nil, fmt.Errorf("element does not exist: [%s]", id)
}

// exists checks if a template with a given ID already exists
func exists(id string) bool {
	for _, t := range elements {
		if t.ID == id {
			return true
		}
	}
	return false
}

func tplId(moduleID, tplID string) string {
	return strings.Join([]string{moduleID, tplID}, "/")
}
