package components

import (
	"fmt"
	"orbital/dashboard/wasm/dom"
	"orbital/dashboard/wasm/events"
	"orbital/dashboard/wasm/storage"
	"syscall/js"
)

type LoginComponentDI struct {
	Events  *events.Event
	Storage storage.Storage
}

type LoginComponent struct {
	secretKey string
	tplDir    string
	events    *events.Event
	store     storage.Storage
}

func (c *LoginComponent) registerEvents() {
	c.events.On("login.show", c.Show)
}

func (c *LoginComponent) Show() {

	fmt.Println("Login Component Show")

	tplObj, err := dom.GetElement(c.tplDir, "main/default")
	if err != nil {
		dom.PrintToConsole(fmt.Sprintf("Error loading template: %s", err))
		return
	}
	htmlEl := tplObj.CloneFromTemplate()

	// Render the template
	renderedElement, err := dom.RenderStatic(htmlEl.Obj, nil)
	if err != nil {
		fmt.Printf("Error rendering login template: %v\n", err)
		return
	}

	loginBtn := renderedElement.Call("querySelector", "#login-button")
	loginBtn.Call("addEventListener", "click", js.FuncOf(c.uiLoginAction))

	c.events.Emit("app.render", renderedElement)
}

func (c *LoginComponent) uiLoginAction(this js.Value, args []js.Value) interface{} {
	input := dom.DocQuerySelectorValue("#public-key", "value")
	fmt.Println("Validate:", input)
	return nil
}

func NewLoginComponents(di LoginComponentDI) {
	c := &LoginComponent{
		events: di.Events,
		store:  di.Storage,
		tplDir: "login",
	}

	c.registerEvents()
}
