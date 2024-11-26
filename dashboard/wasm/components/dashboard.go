package components

import (
	"fmt"
	"orbital/dashboard/wasm/dom"
	"orbital/dashboard/wasm/events"
	"orbital/dashboard/wasm/storage"
)

type DashboardComponentDI struct {
	Events  *events.Event
	Storage storage.Storage
}

type DashboardComponent struct {
	tplDir string
	events *events.Event
	store  storage.Storage
}

func (c *DashboardComponent) registerEvents() {
	c.events.On("dashboard.show", c.Show)
}

func (c *DashboardComponent) Show() {
	tpl, err := dom.GetElement(c.tplDir, "main")
	if err != nil {
		fmt.Println(err)
		return
	}

	htmlEl := tpl.CloneFromTemplate()

	// Render the template
	renderedElement, err := dom.RenderStatic(htmlEl.Obj, nil)
	if err != nil {
		fmt.Printf("Error rendering login template: %v\n", err)
		return
	}

	c.events.Emit("app.render", renderedElement)
}

func NewDashboardComponent(di DashboardComponentDI) {
	c := &DashboardComponent{
		events: di.Events,
		store:  di.Storage,
		tplDir: "dashboard",
	}

	c.registerEvents()
}
