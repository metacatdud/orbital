package main

import (
	"embed"
	"fmt"
	"io/fs"
	"orbital/web/wasm/app"
	"orbital/web/wasm/components"
	"orbital/web/wasm/dom"
	"orbital/web/wasm/events"
	"orbital/web/wasm/storage"
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
			var htmlBin []byte
			htmlBin, err = templateFS.ReadFile(path)
			if err != nil {
				fmt.Printf("Error reading template %s: %v\n", path, err)
				return err
			}

			tmplName := strings.TrimPrefix(path, "templates/")
			tmplName = strings.TrimSuffix(tmplName, ".html")

			err = dom.RegisterElement(tmplName, htmlBin)
			if err != nil {
				fmt.Printf("Error adding template %s: %v\n", path, err)
			}

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
func bootstrapApp(_ js.Value, _ []js.Value) interface{} {

	event := events.New()
	store := storage.NewLocalStorage()
	ws := app.NewWsConn(true)

	orbital := app.NewApp(app.AppDI{
		Events:  event,
		Storage: store,
		WsConn:  ws,
	})

	components.NewDashboardComponent(components.DashboardComponentDI{
		Events:  event,
		Storage: store,
		WsConn:  ws,
	})

	components.NewLoginComponents(components.LoginComponentDI{
		Events:  event,
		Storage: store,
	})

	orbital.Boot()

	return nil
}
