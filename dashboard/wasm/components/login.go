package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Login struct {
	app.Compo
	secretKey string
}

func (l *Login) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Login"),
		app.Input().
			Type("password").
			Value(l.secretKey).
			OnChange(l.ValueTo(&l.secretKey)),
		app.Button().
			Text("Login").
			OnClick(l.onLogin),
	)
}

func (l *Login) onLogin(ctx app.Context, e app.Event) {
	fmt.Println("Clicked login")
}
