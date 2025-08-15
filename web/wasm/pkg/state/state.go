package state

import (
	"orbital/web/wasm/pkg/dom"
	"reflect"
	"sync"
)

type stateItem struct {
	value    any
	oldValue any
	typeRef  reflect.Type
}

type State struct {
	states   map[string]stateItem
	watchers map[string][]func(oldValue, newValue any)
	mu       sync.RWMutex
}

func New() *State {
	return &State{
		states:   make(map[string]stateItem),
		watchers: make(map[string][]func(oldValue, newValue any)),
	}
}

func (s *State) Get(key string) any {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.states[key]
	if !ok {
		return nil
	}
	return item.value
}

func (s *State) GetAll() map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stateCopy := make(map[string]any, len(s.states))
	for k, item := range s.states {
		stateCopy[k] = item.value
	}
	return stateCopy
}

func (s *State) Set(key string, value any) {
	s.mu.Lock()

	newType := reflect.TypeOf(value)
	if newType == nil {
		dom.ConsoleError("nil value provided", key)
		s.mu.Unlock()
		return
	}

	// Reject pointer values.
	if newType.Kind() == reflect.Ptr {
		dom.ConsoleError("pointer types are not supported for key", key, " Received pointer", newType)
		s.mu.Unlock()
		return
	}

	oldItem, exists := s.states[key]

	if exists {
		if oldItem.typeRef != newType {
			dom.ConsoleError("types mismatch for existing state", key, "Received", newType.String(), "Expected", oldItem.typeRef.String())
			s.mu.Unlock()
			return
		}

		// Skip update if the values are deeply equal.
		if reflect.DeepEqual(oldItem.value, value) {
			s.mu.Unlock()
			return
		}
	}

	s.states[key] = stateItem{
		value:    value,
		oldValue: oldItem.value,
		typeRef:  newType,
	}

	watchers := s.watchers[key]
	s.mu.Unlock()

	for _, cb := range watchers {
		go cb(oldItem.value, value)
	}

	// If the new value is a struct check its fields.
	if newType.Kind() == reflect.Struct {
		if exists && oldItem.typeRef.Kind() == reflect.Struct {
			s.setStructObserver(key, oldItem.value, value)
		}
	}

}

func (s *State) Watch(key string, callback func(oldValue, newValue any)) func() {
	s.mu.Lock()
	s.watchers[key] = append(s.watchers[key], callback)
	s.mu.Unlock()

	return func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		list := s.watchers[key]
		for i, cb := range list {
			if reflect.ValueOf(cb).Pointer() == reflect.ValueOf(callback).Pointer() {
				s.watchers[key] = append(list[:i], list[i+1:]...)
				break
			}
		}
	}
}

func (s *State) setStructObserver(key string, oldValue, newValue any) {

	oldVal := reflect.ValueOf(oldValue)
	newVal := reflect.ValueOf(newValue)

	// Ensure both values are structs.
	if oldVal.Kind() != reflect.Struct || newVal.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < oldVal.NumField(); i++ {
		fieldName := oldVal.Type().Field(i).Name
		oldFieldValue := oldVal.Field(i).Interface()
		newFieldValue := newVal.Field(i).Interface()

		// If changes, notify watchers list
		if !reflect.DeepEqual(oldFieldValue, newFieldValue) {
			watchKey := key + "." + fieldName

			s.mu.RLock()
			watchers := s.watchers[watchKey]
			s.mu.RUnlock()

			for _, cb := range watchers {
				go cb(oldFieldValue, newFieldValue)
			}
		}
	}
}

// MergeStateWithData will merge data from state with other data
// the data from State will be overwritten
func MergeStateWithData(stateData map[string]any, data ...map[string]any) map[string]any {
	mergedData := make(map[string]any, len(stateData))
	for k, v := range stateData {
		mergedData[k] = v
	}

	if len(data) > 0 {
		for k, v := range data[0] {
			if _, found := stateData[k]; !found {
				mergedData[k] = v
			}
		}

	}

	return mergedData
}
