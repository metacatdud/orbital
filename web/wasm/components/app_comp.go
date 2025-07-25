package components

import (
	"orbital/web/wasm/orbital"
	"orbital/web/wasm/service"
)

type AppComponent struct {
	*BaseComponent
	id RegKey
}

// TODO: Need to fir a way to add component logic from both:
// - local (embed.FS)
// - user defined (user storage)

func NewAppComponent(di *orbital.Dependency, id RegKey, data service.App) *AppComponent {
	base := NewBaseComponent(di, id, data.Namespace)
	comp := &AppComponent{
		BaseComponent: base,
	}

	return comp
}
