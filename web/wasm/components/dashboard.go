package components

import (
	"fmt"
	"orbital/web/wasm/api"
	"orbital/web/wasm/dom"
	"orbital/web/wasm/domain"
	"orbital/web/wasm/events"
	"syscall/js"
)

type DashboardComponentDI struct {
	Events   *events.Event
	WsConn   *api.WsConn
	AuthRepo domain.AuthRepository
	UserRepo domain.UserRepository
}

type DashboardComponent struct {
	tplDir   string
	events   *events.Event
	wsConn   *api.WsConn
	authRepo domain.AuthRepository
	userRepo domain.UserRepository
}

func (c *DashboardComponent) registerEvents() {
	c.events.On("dashboard.show", c.Show)

	c.wsConn.On("orbital.machine", func(data []byte) {
		dom.PrintToConsole("machine stats", string(data))
	})
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
	//wsTestBtn := dom.ElementQuerySelect(renderedElement, "[data-action='wsTest']")
	//wsTestBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	//	fmt.Println("---")
	//	var authData map[string]string
	//	if err = c.storage.Get("auth", &authData); err != nil {
	//		dom.PrintToConsole("Failed to get public key")
	//		return nil
	//	}
	//
	//	publicKey, _ := cryptographer.NewPublicKeyFromString(authData["publicKey"])
	//	fmt.Printf("publicKey (from string): %+v\n", publicKey.String())
	//	//
	//	msg := app.NewTopicMessage("dashboard.allData", []byte(`req.data`))
	//	msg.PublicKey = publicKey.Compress()
	//
	//	c.wsConn.Send(*msg)
	//	return nil
	//}))

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
		events:   di.Events,
		tplDir:   "dashboard",
		wsConn:   di.WsConn,
		authRepo: di.AuthRepo,
		userRepo: di.UserRepo,
	}

	c.registerEvents()
}
