package components

import "orbital/web/wasm/orbital"

const (
	AboutComponentRegKey RegKey = "aboutComponent"
)

type AboutComponent struct {
	*BaseComponent
}

var _ MetaProvider = (*AboutComponent)(nil)

func NewAboutComponent(di *orbital.Dependency) *AboutComponent {
	base := NewBaseComponent(di, AboutComponentRegKey, "orbital/about/about")
	comp := &AboutComponent{
		BaseComponent: base,
	}

	return comp
}

func ID() RegKey {
	return AboutComponentRegKey
}

func (a *AboutComponent) Title() string {
	return "About"
}

func (a *AboutComponent) Icon() string {
	return "fa-circle-info"
}
