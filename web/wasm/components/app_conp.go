package components

import "orbital/web/wasm/orbital"

type AppComponent struct {
	*BaseComponent

	titleText string
	iconClass string
}

var _ MetaProvider = (*AppComponent)(nil)

func NewAppComponent(di *orbital.Dependency, name RegKey) *AppComponent {
	base := NewBaseComponent(di, name, "dashboard/app/appThumb")
	comp := &AppComponent{
		BaseComponent: base,
	}

	return comp
}

func (a AppComponent) Title() string {
	return a.titleText
}

func (a AppComponent) Icon() string {
	return a.iconClass
}
