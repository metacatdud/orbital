package orbital

import (
	"fmt"
	"orbital/web/wasm/pkg/events"
	"orbital/web/wasm/pkg/state"
	"orbital/web/wasm/pkg/storage"
	"orbital/web/wasm/pkg/transport"
	"orbital/web/wasm/templates"
)

type ServiceName string

type Option func(*Dependency) error

func WithEvents(e *events.Event) Option {
	return func(d *Dependency) error {
		d.Events = e
		return nil
	}
}

func WithState(s *state.State) Option {
	return func(d *Dependency) error {
		d.State = s
		return nil
	}
}

func WithStorage(s storage.Storage) Option {
	return func(d *Dependency) error {
		d.Storage = s
		return nil
	}
}

func WithTemplates(t *templates.Registry) Option {
	return func(d *Dependency) error {
		d.Templates = t
		return nil
	}
}

func WithWs(ws *transport.WsConn) Option {
	return func(d *Dependency) error {
		d.Ws = ws
		return nil
	}
}

func WithService(name ServiceName, instance any) Option {
	return func(d *Dependency) error {
		if instance == nil {
			return fmt.Errorf("service instance must not be nil")
		}

		if _, exists := d.services[name]; exists {
			return fmt.Errorf("service instance already exists")
		}

		d.services[name] = instance

		if sub, ok := instance.(events.EventSubscriber); ok {
			sub.HookEvents(d.Events)

		}
		return nil
	}
}

type Dependency struct {
	Events    *events.Event
	State     *state.State
	Storage   storage.Storage
	Templates *templates.Registry
	Ws        *transport.WsConn
	services  map[ServiceName]any
}

func NewDependency(opts ...Option) (*Dependency, error) {
	tplRegistry, err := templates.NewRegistry()
	if err != nil {
		return nil, err
	}

	d := &Dependency{
		Events:    events.New(),
		State:     state.New(),
		Storage:   storage.NewLocalStorage(),
		Templates: tplRegistry,
		Ws:        transport.NewWsConn(true),
		services:  make(map[ServiceName]any),
	}

	for _, opt := range opts {
		if err = opt(d); err != nil {
			return nil, err
		}
	}

	return d, nil
}

func (d *Dependency) RegisterService(name ServiceName, instance any) error {
	if instance == nil {
		return fmt.Errorf("service instance must not be nil")
	}

	if _, exists := d.services[name]; exists {
		return fmt.Errorf("service instance already exists")
	}

	d.services[name] = instance

	if sub, ok := instance.(events.EventSubscriber); ok {
		sub.HookEvents(d.Events)

	}

	return nil
}

func (d *Dependency) GetService(name ServiceName) (any, error) {
	service, exists := d.services[name]
	if !exists {
		return nil, fmt.Errorf("service does not exist")
	}

	return service, nil
}

func MustGetService[T any](d *Dependency, name ServiceName) T {
	raw, err := d.GetService(name)
	if err != nil {
		panic(err)
	}
	svc, ok := raw.(T)
	if !ok {
		var zero T
		panic(fmt.Sprintf("service %q: got %T, wanted %T", name, raw, zero))
	}
	return svc
}
