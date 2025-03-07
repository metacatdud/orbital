package events

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
)

var Global = New()

type Handler struct {
	id       uint64
	Callback interface{}
	Once     bool
}

type Event struct {
	mu        sync.RWMutex
	listeners map[string][]*Handler
	counter   uint64
}

func New() *Event {
	return &Event{
		listeners: make(map[string][]*Handler),
	}
}

// On creates a listener for an event and returns an unsubscribe function.
func (e *Event) On(eventName string, handler interface{}) func() {
	return e.register(eventName, handler, false)
}

// Once creates a one-time listener for an event and returns an unsubscribe function.
func (e *Event) Once(eventName string, handler interface{}) func() {
	return e.register(eventName, handler, true)
}

// Remove deletes all listeners for a given event.
func (e *Event) Remove(eventName string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.listeners, eventName)
}

func (e *Event) Emit(eventName string, params ...interface{}) {

	e.mu.RLock()

	handlers, ok := e.listeners[eventName]
	if !ok {
		e.mu.RUnlock()
		return
	}

	// Copy the slice to avoid race conditions if it gets modified during iteration.
	copiedHandlers := make([]*Handler, len(handlers))
	copy(copiedHandlers, handlers)

	e.mu.RUnlock()

	var onceIDs []uint64
	for _, h := range copiedHandlers {
		e.callHandler(h.Callback, params)
		if h.Once {
			onceIDs = append(onceIDs, h.id)
		}
	}

	if len(onceIDs) > 0 {
		e.mu.Lock()
		currentHandlers, ok := e.listeners[eventName]
		if ok {
			newHandlers := currentHandlers[:0]
			for _, h := range currentHandlers {
				remove := false
				for _, id := range onceIDs {
					if h.id == id {
						remove = true
						break
					}
				}
				if !remove {
					newHandlers = append(newHandlers, h)
				}
			}

			e.listeners[eventName] = newHandlers
		}
		e.mu.Unlock()
	}
}

func (e *Event) callHandler(handler interface{}, params []interface{}) {
	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()

	if len(params) != handlerType.NumIn() {
		panic(fmt.Sprintf("handler expects %d parameters, got %d", handlerType.NumIn(), len(params)))
	}

	args := make([]reflect.Value, len(params))
	for i, p := range params {
		args[i] = reflect.ValueOf(p)
	}

	handlerValue.Call(args)
}

func (e *Event) register(eventName string, handler interface{}, once bool) func() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if reflect.TypeOf(handler).Kind() != reflect.Func {
		panic(fmt.Sprintf("handler for event %s is not a function", eventName))
	}

	id := atomic.AddUint64(&e.counter, 1)
	h := &Handler{
		id:       id,
		Callback: handler,
		Once:     once,
	}

	e.listeners[eventName] = append(e.listeners[eventName], h)

	return func() {
		e.off(eventName, id)
	}
}

func (e *Event) off(eventName string, handlerID uint64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	handlers, ok := e.listeners[eventName]
	if !ok {
		return
	}

	// Filter out the handler with the matching ID.
	newHandlers := handlers[:0]
	for _, h := range handlers {
		if h.id != handlerID {
			newHandlers = append(newHandlers, h)
		}
	}
	e.listeners[eventName] = newHandlers
}
