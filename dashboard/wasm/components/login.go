package components

import (
	"encoding/json"
	"fmt"
	"orbital/dashboard/wasm/api"
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
	async     *api.Async
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

	req := map[string]string{
		"publicKey": input.String(),
	}

	reqBin, err := json.Marshal(req)
	if err != nil {
		dom.PrintToConsole(fmt.Sprintf("Error serializing public key: %v", err))
		return nil
	}

	resultChan := c.async.RunWithResult(func() (interface{}, error) {
		client := api.NewAPI("/rpc/AuthService/Auth")
		res, err := client.Do(reqBin, nil)
		if err != nil {
			dom.PrintToConsole(fmt.Sprintf("Error calling HelloService: %v", err))
			return nil, err
		}

		resMap := make(map[string]interface{})
		_ = json.Unmarshal(res, &resMap)

		fmt.Println("Response:", resMap)
		return resMap, nil
	})

	go func() {
		result := <-resultChan
		if result.Err != nil {
			fmt.Printf("Error validating public key: %v\n", result.Err)
			return
		}

		fmt.Printf("Validation successful, response: %s\n", result.Value)
	}()

	return nil
}

func NewLoginComponents(di LoginComponentDI) {
	c := &LoginComponent{
		async:  api.NewAsync(),
		events: di.Events,
		store:  di.Storage,
		tplDir: "login",
	}

	c.registerEvents()
}
