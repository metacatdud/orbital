package component

import (
	"syscall/js"
)

// Component interface for basic components
type Component interface {
	ID() string
	Mount(container *js.Value) error
	Unmount() error
	Namespace() string
	Render() error
}

// StateControl implementation for components with state manager
// Should implement `state.State`
type StateControl interface {
	BindStateWatch()
	UnbindStateWatch()
}

// EventControl implementation for components with events
// Should implement `events.Event`
type EventControl interface {
	BindEvents()
	UnbindEvents()
}
