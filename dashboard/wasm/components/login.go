package components

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Login struct {
	app.Compo
	secretKey string
}

func (comp *Login) OnMount(ctx *app.Context) {

}

func (comp *Login) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Login"),
		app.Input().
			Type("password").
			Value(comp.secretKey).
			OnChange(comp.ValueTo(&comp.secretKey)),
		app.Button().
			Text("Login").
			OnClick(comp.onLogin),
	)
}

func (comp *Login) onLogin(ctx app.Context, e app.Event) {
	fmt.Println("Clicked login 123")
	ctx.NewActionWithValue("dummyHandler", "login")
	fmt.Println("after dummy")
}
