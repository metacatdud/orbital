package component

import (
	"syscall/js"
)

// Component interface for basic components
type Component interface {
	Mount(container *js.Value) error
	Unmount() error
	Namespace() string
	Render(data ...map[string]interface{}) (string, error)
}

// ContainerComponent implementation for components with children
type ContainerComponent interface {
	Component
	AddChild(child Component) error
	RemoveChild(child Component) error
	Children() []Component
}

// StaticComponent implementation for simple static components
type StaticComponent interface {
	Component
}

// StatefulComponent implementation for components with state manager
type StatefulComponent interface {
	Component
	RegisterStateWatch()
}

// EventComponent implementation for components with events
type EventComponent interface {
	Component
	RegisterEvents()
}
