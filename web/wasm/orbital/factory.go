package orbital

import (
	"fmt"
)

type FactoryFn func(di *Dependency, params ...interface{}) Component

type Factory struct {
	di       *Dependency
	registry map[string]FactoryFn
}

func NewFactory(di *Dependency) *Factory {
	return &Factory{
		di:       di,
		registry: make(map[string]FactoryFn),
	}
}

func (f *Factory) Create(name string, params ...interface{}) (Component, error) {
	create, exists := f.registry[name]
	if !exists {
		return nil, fmt.Errorf("component '%s' is not registered", name)
	}

	return create(f.di, params...), nil
}

func (f *Factory) Register(name string, component FactoryFn) {
	f.registry[name] = component
}
