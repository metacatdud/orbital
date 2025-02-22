package events

import (
	"fmt"
	"reflect"
	"sync"
)

var Global = New()

type Handler struct {
	Callback interface{}
	Once     bool
}

type Event struct {
	mu        sync.RWMutex
	listeners map[string][]*Handler
}

func New() *Event {
	return &Event{
		listeners: make(map[string][]*Handler),
	}
}

// On create a listener for an event
func (e *Event) On(eventName string, handler interface{}) {
	e.register(eventName, handler, false)
}

// Once create a listener for an event which will be used once
func (e *Event) Once(eventName string, handler interface{}) {
	e.register(eventName, handler, true)
}

// Off remove a listen from event stack
func (e *Event) Off(eventName string, handler interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if handlers, ok := e.listeners[eventName]; ok {
		newHandlers := make([]*Handler, 0)

		for _, h := range handlers {
			if h.Callback != handler {
				newHandlers = append(newHandlers, h)
			}
		}

		e.listeners[eventName] = newHandlers
	}
}

// Remove an event from listeners stack with all of its handlers
func (e *Event) Remove(eventName string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.listeners, eventName)
}

func (e *Event) Emit(eventName string, params ...interface{}) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if handlers, ok := e.listeners[eventName]; ok {
		eventsToRemove := []*Handler{}

		for _, h := range handlers {
			e.callHandler(h.Callback, params)
			if h.Once {
				eventsToRemove = append(eventsToRemove, h)
			}
		}

		if len(eventsToRemove) > 0 {
			e.mu.RUnlock() //Unlock read lock to switch to write locker
			e.mu.Lock()

			eventsRemaining := make([]*Handler, 0)
			for _, h := range handlers {
				rm := false

				for _, hRem := range eventsToRemove {
					if hRem == h {
						rm = true
						break
					}
				}

				if !rm {
					eventsRemaining = append(eventsRemaining, h)
				}
			}

			e.listeners[eventName] = eventsRemaining
			e.mu.Unlock()
			e.mu.RLock() // Switch locker back to read lock
		}
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

func (e *Event) register(eventName string, handler interface{}, once bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if reflect.TypeOf(handler).Kind() != reflect.Func {
		panic(fmt.Sprintf("handler for event %s is not a function", eventName))
	}

	e.listeners[eventName] = append(e.listeners[eventName], &Handler{
		Callback: handler,
		Once:     once,
	})
}
