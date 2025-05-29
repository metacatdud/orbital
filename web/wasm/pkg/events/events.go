package events

import (
	"fmt"
	"orbital/web/wasm/pkg/dom"
	"reflect"
	"sync"
	"sync/atomic"
)

type EventSubscriber interface {
	HookEvents(e *Event)
}

type Handler struct {
	id       uint64
	Callback any
	Once     bool
}

type Event struct {
	mu        sync.RWMutex
	listeners map[string][]*Handler
	counter   atomic.Uint64
}

func New() *Event {
	return &Event{
		listeners: make(map[string][]*Handler),
	}
}

// On creates a listener for an event and returns an unsubscribe function.
func (e *Event) On(eventName string, handler any) func() {
	return e.register(eventName, handler, false)
}

// Once creates a one-time listener for an event and returns an unsubscribe function.
func (e *Event) Once(eventName string, handler any) func() {
	return e.register(eventName, handler, true)
}

// RemoveAll deletes all listeners for a given event.
func (e *Event) RemoveAll(eventName string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.listeners, eventName)
}

func (e *Event) Emit(eventName string, params ...any) {

	e.mu.RLock()

	handlers, ok := e.listeners[eventName]
	if !ok {
		e.mu.RUnlock()
		return
	}

	// Copy the slice to avoid race conditions if handler gets
	// modified during iteration.
	copiedHandlers := make([]*Handler, len(handlers))
	copy(copiedHandlers, handlers)

	e.mu.RUnlock()

	var oneTimeEventIDs []uint64
	for _, h := range copiedHandlers {
		e.callHandler(h.Callback, params)
		if h.Once {
			oneTimeEventIDs = append(oneTimeEventIDs, h.id)
		}
	}

	if len(oneTimeEventIDs) > 0 {
		e.mu.Lock()
		defer e.mu.Unlock()

		currentHandlers := e.listeners[eventName]
		newHandlers := currentHandlers[:0]
		for _, h := range currentHandlers {
			keep := true
			for _, id := range oneTimeEventIDs {
				if id == h.id {
					keep = false
					break
				}
			}

			if keep {
				newHandlers = append(newHandlers, h)
			}
		}

		e.listeners[eventName] = newHandlers
	}
}

func (e *Event) callHandler(handler any, params []any) {
	hVal := reflect.ValueOf(handler)
	hType := hVal.Type()

	if len(params) != hType.NumIn() {
		dom.ConsoleWarn(fmt.Sprintf("handler expects %d parameters, got %d", hType.NumIn(), len(params)))
		return
	}

	args := make([]reflect.Value, len(params))
	for i, p := range params {
		args[i] = reflect.ValueOf(p)
	}

	hVal.Call(args)
}

func (e *Event) register(eventName string, handler any, once bool) func() {
	if reflect.TypeOf(handler).Kind() != reflect.Func {
		dom.ConsoleWarn(fmt.Sprintf("handler for event %s is not a function", eventName))
		return func() {}
	}

	id := e.counter.Add(1)
	h := &Handler{
		id:       id,
		Callback: handler,
		Once:     once,
	}

	e.mu.Lock()
	e.listeners[eventName] = append(e.listeners[eventName], h)
	e.mu.Unlock()

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

	newHandlers := handlers[:0]
	for _, h := range handlers {
		if h.id != handlerID {
			newHandlers = append(newHandlers, h)
		}
	}
	e.listeners[eventName] = newHandlers
}
