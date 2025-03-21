package orbital

import (
	"fmt"
	"sync"
)

type Mod interface {
	ID() string
}

type ModFactoryFn func(di *Dependency, params ...interface{}) Mod

type registry struct {
	di        *Dependency
	factories map[string]ModFactoryFn
	mu        sync.RWMutex
}

var globalReg *registry

func NewRegistry(di *Dependency) {
	fmt.Println("Init Registry")
	globalReg = &registry{
		di:        di,
		factories: make(map[string]ModFactoryFn),
	}
}

func Register(moduleID string, factory ModFactoryFn) error {
	if moduleID == "" {
		return fmt.Errorf("module id cannot be empty")
	}

	globalReg.mu.Lock()
	defer globalReg.mu.Unlock()

	if _, exists := globalReg.factories[moduleID]; exists {
		return fmt.Errorf("[%w]: %s", ErrRegDuplicateID, moduleID)
	}

	globalReg.factories[moduleID] = factory
	return nil
}

func Lookup[T Mod](modID string, params ...interface{}) (T, error) {
	globalReg.mu.RLock()
	factory, exists := globalReg.factories[modID]
	globalReg.mu.RUnlock()

	if !exists {
		var zero T
		return zero, fmt.Errorf("[%w]: %s", ErrRegNotFound, modID)
	}

	mod := factory(globalReg.di, params...)
	trueMod, ok := mod.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("[%w]: %s", ErrRegWrongType, modID)
	}

	return trueMod, nil
}

//type FactoryFn func(di *Dependency, params ...interface{}) Component
//
//type Factory struct {
//	di       *Dependency
//	registry map[string]FactoryFn
//}
//
//func NewFactory(di *Dependency) *Factory {
//	return &Factory{
//		di:       di,
//		registry: make(map[string]FactoryFn),
//	}
//}
//
//func (f *Factory) Create(name string, params ...interface{}) (Component, error) {
//	create, exists := f.registry[name]
//	if !exists {
//		return nil, fmt.Errorf("component '%s' is not registered", name)
//	}
//
//	return create(f.di, params...), nil
//}
//
//func (f *Factory) Register(name string, component FactoryFn) {
//	f.registry[name] = component
//}
