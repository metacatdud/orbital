package main

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"orbital/dashboard/wasm/components"
)

func main() {
	app.Route("/", func() app.Composer {
		return &components.Login{}
	})
	//app.Route("/dashboard", &components.Dashboard{})
	//app.Route("/settings", &components.Settings{})

	app.RunWhenOnBrowser()
}
