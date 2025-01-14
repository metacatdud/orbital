package components

import (
	"fmt"
	"orbital/pkg/cryptographer"
	"orbital/web/wasm/app"
	"orbital/web/wasm/dom"
	"orbital/web/wasm/events"
	"orbital/web/wasm/storage"
	"syscall/js"
)

type DashboardComponentDI struct {
	Events  *events.Event
	Storage storage.Storage
	WsConn  *app.WsConn
}

type DashboardComponent struct {
	tplDir  string
	events  *events.Event
	storage storage.Storage
	wsConn  *app.WsConn
}

func (c *DashboardComponent) registerEvents() {
	c.events.On("dashboard.show", c.Show)
}

func (c *DashboardComponent) Show() {
	tpl, err := dom.GetElement(c.tplDir, "main/default")
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

	// Register UI events
	toggleProfile := dom.ElementQuerySelect(renderedElement, "[data-action='toggleProfile']")
	toggleProfile.Call("addEventListener", "click", js.FuncOf(c.uiToggleProfile))

	avatarCloseOverlay := dom.ElementQuerySelect(renderedElement, "[data-id='avatarCloseOverlay']")
	avatarCloseOverlay.Call("addEventListener", "click", js.FuncOf(c.uiToggleProfileClose))

	avatarLogoutBtn := dom.ElementQuerySelect(renderedElement, "[data-action='logout']")
	avatarLogoutBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		c.events.Emit("user.logout")
		return nil
	}))

	// test for websocket call
	wsTestBtn := dom.ElementQuerySelect(renderedElement, "[data-action='wsTest']")
	wsTestBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("---")
		var authData map[string]string
		if err = c.storage.Get("auth", &authData); err != nil {
			dom.PrintToConsole("Failed to get public key")
			return nil
		}

		publicKey, _ := cryptographer.NewPublicKeyFromString(authData["publicKey"])
		fmt.Printf("publicKey (from string): %+v\n", publicKey.String())
		//
		msg := app.NewTopicMessage("dashboard.allData", []byte(`req.data`))
		msg.PublicKey = publicKey.Compress()

		c.wsConn.Send(*msg)
		return nil
	}))

	c.events.Emit("app.render", renderedElement)
}

func (c *DashboardComponent) uiToggleProfile(this js.Value, args []js.Value) interface{} {

	dom.Show("[data-id='avatarMenu']")
	dom.Show("[data-id='avatarCloseOverlay']")

	return nil
}

func (c *DashboardComponent) uiToggleProfileClose(this js.Value, args []js.Value) interface{} {
	dom.Hide("[data-id='avatarMenu']")
	dom.Hide("[data-id='avatarCloseOverlay']")

	//avatarCloseOverlay := dom.DocQuerySelector("[data-id='avatarCloseOverlay']")
	//avatarCloseOverlay.Call("removeEventListener", "click", js.FuncOf(c.uiToggleProfileClose))

	return nil
}

func NewDashboardComponent(di DashboardComponentDI) {
	c := &DashboardComponent{
		events:  di.Events,
		storage: di.Storage,
		tplDir:  "dashboard",
		wsConn:  di.WsConn,
	}

	c.registerEvents()
}
