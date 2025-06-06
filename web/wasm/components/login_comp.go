package components

import (
	"bytes"
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/dom"
	"orbital/web/wasm/pkg/transport"
	"orbital/web/wasm/service"
	"syscall/js"
)

const (
	LoginComponentRegKey RegKey = "loginComponent"
)

type LoginComponent struct {
	*BaseComponent
	authSvc *service.AuthService
}

var _ MetaProvider = (*LoginComponent)(nil)

func NewLoginComponent(di *orbital.Dependency) *LoginComponent {
	base := NewBaseComponent(di, LoginComponentRegKey, "auth/auth/default")

	authSvc := orbital.MustGetService[*service.AuthService](di, service.AuthServiceKey)

	comp := &LoginComponent{
		BaseComponent: base,
		authSvc:       authSvc,
	}

	comp.bindUIEvents()
	return comp
}

func (comp *LoginComponent) ID() RegKey {
	return LoginComponentRegKey
}

func (comp *LoginComponent) Title() string { return "Login" }

func (comp *LoginComponent) Icon() string { return "fa-lock" }

func (comp *LoginComponent) Unmount() error {
	return comp.BaseComponent.Unmount()
}

func (comp *LoginComponent) renderError(errType, msg string) {
	tpl, err := comp.DI.Templates.Get("auth/auth/errorMsg")
	if err != nil {
		dom.ConsoleError("cannot load template", err.Error())
		return
	}

	var buf bytes.Buffer
	data := map[string]interface{}{"type": errType, "message": msg}
	if err = tpl.Execute(&buf, data); err != nil {
		dom.ConsoleError("cannot execute template", err.Error())
		return
	}

	container := comp.GetContainer("errorMessage")
	if container.IsNull() {
		dom.ConsoleError("cannot find errorMessage container")
		return
	}

	dom.SetInnerHTML(container, buf.String())
	dom.RemoveClass(container, "hide")

}

func (comp *LoginComponent) clearError() {
	container := comp.GetContainer("errorMessage")
	if container.IsNull() {
		dom.ConsoleError("cannot find errorMessage container")
		return
	}

	dom.AddClass(container, "hide")
	dom.SetInnerHTML(container, "")
}

func (comp *LoginComponent) bindUIEvents() {
	comp.AddEventHandler("[data-action='login']", "click", comp.uiEventLogin)
	comp.AddEventHandler("[data-action='about']", "click", comp.uiEventAbout)
}

func (comp *LoginComponent) uiEventLogin(_ js.Value, args []js.Value) interface{} {
	var async transport.Async
	async.Async(func() {

		skInput := dom.GetValue("input", "privateKey")
		res, err := comp.authSvc.Login(service.LoginReq{
			SecretKey: skInput,
		})

		if err != nil {
			comp.renderError("auth.failed", err.Error())
			return
		}

		if res.Error != nil {
			comp.renderError(res.Error.Type, res.Error.Msg)
			return
		}

		comp.clearError()
		comp.DI.State.Set("state:isAuthenticated", true)

		return
	})
	return nil
}

func (comp *LoginComponent) uiEventAbout(_ js.Value, args []js.Value) interface{} {
	comp.DI.State.Set("state:overlay:currentChild", AboutComponentRegKey)
	return nil
}
