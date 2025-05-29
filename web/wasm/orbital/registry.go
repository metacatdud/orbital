package orbital

// TODO: Remove this. Not used anymore
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

// TODO: Add multiple properties for each mod type this should simplify the retrieval.
// might not work as each mod has it's own functionalities
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

	var zero T
	if !exists {
		return zero, fmt.Errorf("[%w]: %s", ErrRegNotFound, modID)
	}

	mod := factory(globalReg.di, params...)
	trueMod, ok := mod.(T)
	if !ok {
		return zero, fmt.Errorf("[%w]: %s", ErrRegWrongType, modID)
	}

	return trueMod, nil
}
