package components

import "orbital/web/wasm/orbital"

const (
	DashboardComponentRegKey RegKey = "dashboardComponent"
)

type DashboardComponent struct {
	*BaseComponent

	// TODO: Add children here
}

func NewDashboardComponent(di *orbital.Dependency) *DashboardComponent {
	base := NewBaseComponent(di, DashboardComponentRegKey, "dashboard/main/default")

	comp := &DashboardComponent{BaseComponent: base}

	return comp
}
