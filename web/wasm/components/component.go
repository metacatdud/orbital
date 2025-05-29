package components

import (
	"syscall/js"
)

type Component interface {
	ID() RegKey
	Mount(container *js.Value) error
	Unmount() error
}

type MetaProvider interface {
	Title() string
	Icon() string
}

type ParentRenderObserver interface {
	OnParentRender()
}
