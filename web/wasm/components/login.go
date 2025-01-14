package components

import (
	"encoding/json"
	"fmt"
	"orbital/web/wasm/api"
	"orbital/web/wasm/dom"
	"orbital/web/wasm/events"
	"orbital/web/wasm/storage"
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
	c.events.On("user.logon", c.UserLogon)
	c.events.On("user.logout", c.UserLogout)
}

func (c *LoginComponent) Show() {
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

	loginBtn := dom.ElementQuerySelect(renderedElement, "[data-action='login']")
	loginBtn.Call("addEventListener", "click", js.FuncOf(c.uiLoginAction))

	c.events.Emit("app.render", renderedElement)
}

func (c *LoginComponent) UserLogon(u *User, errs map[string]string) {
	errPlaceholder := dom.DocQuerySelector("#login-error")
	errPlaceholder.Set("innerHTML", "")

	if u != nil {
		dom.PrintToConsole("Login as:", u.Name, "Access", u.Access)

		if err := c.store.Set("auth", u); err != nil {
			dom.PrintToConsole(fmt.Sprintf("Error storing auth data: %v", err))
			return
		}
		
		c.events.Emit("navigate", "dashboard")
		return
	}

	tplObj, err := dom.GetElement(c.tplDir, "main/error-msg")
	if err != nil {
		dom.PrintToConsole(fmt.Sprintf("Error loading template: %s", err))
		return
	}

	for errTyp, errMsg := range errs {
		htmlEl := tplObj.CloneFromTemplate()

		htmlNode := dom.ElementQuerySelect(htmlEl.Obj, fmt.Sprintf("[data-errTyp='%s']", errTyp))
		htmlNode.Set("textContent", errMsg)

		dom.AppendChild("#login-error", htmlNode)
	}

	dom.Show("#login-error")
}

func (c *LoginComponent) UserLogout() {
	if err := c.store.Del("auth"); err != nil {
		dom.PrintToConsole(fmt.Sprintf("Error deleting auth data: %v", err))
		return
	}

	// Close websocket on user logout
	// TODO: Check if it does make sense as this socket is meant only for one user
	c.events.Emit("ws.close")

	c.Show()
}

func (c *LoginComponent) uiLoginAction(this js.Value, args []js.Value) interface{} {
	input := dom.DocQuerySelectorValue("[data-input='privateKey']", "value")

	req := &LoginReq{
		PublicKey: input.String(),
	}

	reqBin, err := json.Marshal(req)
	if err != nil {
		dom.PrintToConsole(fmt.Sprintf("Error serializing public key: %v", err))
		return nil
	}

	c.async.Run(func() {
		client := api.NewAPI("/rpc/AuthService/Auth")
		var res []byte
		res, err = client.Do(reqBin, nil)
		if err != nil {
			dom.PrintToConsole(fmt.Sprintf("Error calling HelloService: %v", err))
			return
		}

		resMap := &LoginResp{}
		_ = json.Unmarshal(res, resMap)

		c.events.Emit("user.logon", resMap.User, resMap.Err)
		return
	})

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

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
	Access    string `json:"access"`
}

type LoginReq struct {
	PublicKey string `json:"publicKey"`
}

type LoginResp struct {
	User *User             `json:"user"`
	Err  map[string]string `json:"error"`
}
