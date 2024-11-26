package dom

import (
	"fmt"
	"html/template"
	"strings"
	"syscall/js"
)

func RenderStatic(tplObj js.Value, data interface{}) (js.Value, error) {

	if !tplObj.Truthy() {
		return js.Null(), fmt.Errorf("template object is not valid")
	}

	htmlContent := tplObj.Get("outerHTML").String()
	if htmlContent == "" {
		return js.Null(), fmt.Errorf("template content is empty")
	}

	tpl, err := template.New("_temp").Parse(htmlContent)
	if err != nil {
		return js.Null(), fmt.Errorf("failed to parse HTML into template: %v", err)
	}

	var renderedHTML strings.Builder
	err = tpl.Execute(&renderedHTML, data)
	if err != nil {
		return js.Null(), fmt.Errorf("failed to render template: %v", err)
	}

	tmpNode := Document().Obj.Call("createElement", "div")
	tmpNode.Set("innerHTML", renderedHTML.String())

	firstChild := tmpNode.Get("firstChild")
	if !firstChild.Truthy() {
		return js.Null(), fmt.Errorf("failed to create JS element from rendered HTML")
	}

	return firstChild, nil
}
