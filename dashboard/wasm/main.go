package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"orbital/dashboard/wasm/dom"
	"strings"
	"syscall/js"
)

//go:embed templates/*
var templateFS embed.FS

func main() {
	loadTemplates()

	js.Global().Set("bootstrapApp", js.FuncOf(bootstrapApp))

	select {}
}

func loadTemplates() {
	err := fs.WalkDir(templateFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, ".html") {
			var tmpl *template.Template
			tmpl, err = template.ParseFS(templateFS, path)
			if err != nil {
				return err
			}

			// Render the template to a temporary <div> as js.Value
			var buf bytes.Buffer
			err = tmpl.Execute(&buf, nil)
			if err != nil {
				fmt.Printf("Error rendering template %s: %v\n", path, err)
				return err
			}

			// Create a temporary <div> and set its HTML to the rendered content
			//tempContainer := dom.Document().Obj.Call("createElement", "div")
			//tempContainer.Set("innerHTML", buf.String())

			// Register each <template data-template="..."> within this file
			tmplName := strings.TrimPrefix(path, "templates/")
			tmplName = strings.TrimSuffix(tmplName, ".html")

			err = dom.AddModuleTemplate(tmplName, buf.Bytes())
			if err != nil {
				fmt.Printf("Error adding template %s: %v\n", path, err)
			}

			fmt.Printf("Loaded: %s\n", tmplName)
			return nil
		}

		return nil
	})
	if err != nil {
		fmt.Println("Templates load fail:", err.Error())
		return
	}
}

// bootstrapApp load the entrypoint
func bootstrapApp(this js.Value, args []js.Value) interface{} {
	dom.Clear("#app")

	clonedContent, err := dom.GetTemplate("dashboard/main", "custom")
	if err != nil {
		fmt.Println("Cannot load template:", err.Error())
		return nil
	}

	if clonedContent.Truthy() {
		dom.AppendChild("#app", clonedContent)
	}

	dom.Hide("#loading")
	dom.Show("#app")

	return nil
}
