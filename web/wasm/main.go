package main

import (
	"orbital/web/wasm/components"
	"orbital/web/wasm/domain"
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/pkg/dom"
	"syscall/js"
)

func main() {

	js.Global().Set("bootstrapApp", js.FuncOf(bootstrapApp))

	select {}
}

// bootstrapApp load the entrypoint
func bootstrapApp(_ js.Value, _ []js.Value) interface{} {

	// Dependencies
	di, err := buildDependencies()
	if err != nil {
		dom.ConsoleError("Dependency build error", err.Error())
		return nil
	}

	_ = setupFactory(di)

	// Services
	_ = domain.NewAuthService(di)
	_ = domain.NewMachineService(di)

	// Boot Orbital
	app, err := orbital.NewOrbital(di)
	if err != nil {
		dom.ConsoleError("Cannot create orbital app", err.Error())
		return nil
	}

	app.Boot()

	return nil
}

// setupFactory returns a factory for components with dynamic location
func setupFactory(di *orbital.Dependency) *orbital.Factory {

	// Use proper factory instance
	factory := di.Factory()

	factory.Register("loginComp", func(di *orbital.Dependency, params ...interface{}) orbital.Component {
		return components.NewLoginComponent(di)
	})

	return factory
}

func buildDependencies() (*orbital.Dependency, error) {
	return orbital.NewDependencyWithDefaults()
}
