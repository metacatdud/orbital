package components

import (
	"encoding/json"
	"fmt"
	"orbital/pkg/cryptographer"
	"orbital/pkg/proto"
	"orbital/web/wasm/api"
	"orbital/web/wasm/dom"
	"orbital/web/wasm/domain"
	"orbital/web/wasm/events"
	"orbital/web/wasm/storage"
	"syscall/js"
)

type LoginComponentDI struct {
	Events   *events.Event
	Storage  storage.Storage
	AuthRepo domain.AuthRepository
	UserRepo domain.UserRepository
}

type LoginComponent struct {
	tplDir   string
	events   *events.Event
	authRepo domain.AuthRepository
	userRepo domain.UserRepository
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

func (c *LoginComponent) UserLogon(u *User, secretKey string, errs map[string]string) {
	errPlaceholder := dom.DocQuerySelector("#login-error")
	errPlaceholder.Set("innerHTML", "")

	if u != nil {
		dom.PrintToConsole("Login as:", u.Name, "Access", u.Access)

		if err := c.authRepo.Save(domain.Auth{SecretKey: secretKey}); err != nil {
			dom.PrintToConsole(fmt.Sprintf("Error saving login credentials: %v", err))
		}

		domainUser := domain.User{
			ID:        u.ID,
			Name:      u.Name,
			PublicKey: u.PublicKey,
			Access:    u.Access,
		}
		if err := c.userRepo.Save(domainUser); err != nil {
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
	if err := c.authRepo.Remove(); err != nil {
		dom.PrintToConsole(fmt.Sprintf("Error removing auth data: %v", err))
		return
	}

	if err := c.userRepo.Remove(); err != nil {
		dom.PrintToConsole(fmt.Sprintf("Error deleting auth data: %v", err))
		return
	}

	c.Show()
}

func (c *LoginComponent) uiLoginAction(this js.Value, args []js.Value) interface{} {
	input := dom.DocQuerySelectorValue("[data-input='privateKey']", "value")

	secretKey, err := cryptographer.NewPrivateKeyFromString(input.String())
	if err != nil {
		dom.PrintToConsole(err.Error())
		return nil
	}

	loginReq := &LoginReq{
		PublicKey: secretKey.PublicKey().String(),
	}

	loginReqBin, err := json.Marshal(loginReq)
	if err != nil {
		dom.PrintToConsole(fmt.Sprintf("Error serializing public key: %v", err))
		return nil
	}

	req := &proto.Message{
		PublicKey: secretKey.PublicKey().Compress(),
		V:         1,
		Body:      loginReqBin,
		Timestamp: proto.TimestampNow(),
	}

	if err = req.Sign(secretKey.Bytes()); err != nil {
		dom.PrintToConsole(fmt.Sprintf("Error signing login credentials: %v", err))
		return nil
	}

	reqBin, err := json.Marshal(req)
	if err != nil {
		dom.PrintToConsole(fmt.Sprintf("Error serializing login credentials: %v", err))
		return nil
	}

	async := api.NewAsync()
	async.Run(func() {
		client := api.NewAPI("/rpc/AuthService/Auth")
		var res []byte
		res, err = client.Do(reqBin, nil)
		if err != nil {
			dom.PrintToConsole(fmt.Sprintf("Error calling HelloService: %v", err))
			return
		}

		resMap := &LoginResp{}
		_ = json.Unmarshal(res, resMap)

		c.events.Emit(
			"user.logon",
			resMap.User,
			input.String(),
			resMap.Err,
		)
		return
	})

	return nil
}

func NewLoginComponents(di LoginComponentDI) {
	c := &LoginComponent{
		events:   di.Events,
		tplDir:   "login",
		authRepo: di.AuthRepo,
		userRepo: di.UserRepo,
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
