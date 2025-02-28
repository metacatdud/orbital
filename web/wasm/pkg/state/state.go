package state

import (
	"reflect"
	"sync"
)

type stateItem struct {
	value    interface{}
	oldValue interface{}
	typeRef  reflect.Type
}

type State struct {
	states   map[string]stateItem
	watchers map[string][]func(oldValue, newValue interface{})
	mu       sync.RWMutex
}

func New() *State {
	return &State{
		states:   make(map[string]stateItem),
		watchers: make(map[string][]func(oldValue, newValue interface{})),
	}
}

func (s *State) Get(key string) interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.states[key]
	if !ok {
		return nil
	}
	return item.value
}

func (s *State) GetAll() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stateCopy := make(map[string]interface{}, len(s.states))
	for k, item := range s.states {
		stateCopy[k] = item.value
	}
	return stateCopy
}

func (s *State) Set(key string, value interface{}) {
	s.mu.Lock()

	oldItem, exists := s.states[key]
	newType := reflect.TypeOf(value)

	// If no change, skip
	if exists && reflect.DeepEqual(oldItem.value, value) {
		s.mu.Unlock()
		return
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
		s.setStructObserver(key, oldItem.value, value)
	}

}

func (s *State) Watch(key string, callback func(oldValue, newValue interface{})) func() {
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

func (s *State) setStructObserver(key string, oldValue, newValue interface{}) {
	oldVal := reflect.ValueOf(oldValue)
	newVal := reflect.ValueOf(newValue)

	// If pointers, get the element.
	if oldVal.Kind() == reflect.Ptr {
		oldVal = oldVal.Elem()
	}

	if newVal.Kind() == reflect.Ptr {
		newVal = newVal.Elem()
	}

	// Ensure struct
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
