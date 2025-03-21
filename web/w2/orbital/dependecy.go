package orbital

import (
	"orbital/web/w2/pkg/events"
	"orbital/web/w2/pkg/state"
	"orbital/web/w2/pkg/storage"
	"orbital/web/w2/pkg/transport"
	"orbital/web/w2/templates"
)

type Packages struct {
	Events      *events.Event
	State       *state.State
	Storage     storage.Storage
	TplRegistry *templates.Registry
	Ws          *transport.WsConn
}

type Dependency struct {
	events      *events.Event
	state       *state.State
	storage     storage.Storage
	tplRegistry *templates.Registry
	ws          *transport.WsConn
}

func (dep *Dependency) State() *state.State {
	return dep.state
}

func (dep *Dependency) TplRegistry() *templates.Registry {
	return dep.tplRegistry
}

func (dep *Dependency) Storage() storage.Storage {
	return dep.storage
}

func (dep *Dependency) Events() *events.Event {
	return dep.events
}

func (dep *Dependency) Ws() *transport.WsConn {
	return dep.ws
}

func NewDependency(pkgs Packages) (*Dependency, error) {
	d := &Dependency{
		events:      pkgs.Events,
		state:       pkgs.State,
		storage:     pkgs.Storage,
		tplRegistry: pkgs.TplRegistry,
		ws:          pkgs.Ws,
	}

	return d, nil
}

func NewDependencyWithDefaults() (*Dependency, error) {
	tplRegistry, err := templates.NewRegistry()
	if err != nil {
		return nil, err
	}

	d := &Dependency{
		events:      events.New(),
		state:       state.New(),
		storage:     storage.NewLocalStorage(),
		tplRegistry: tplRegistry,
		ws:          transport.NewWsConn(true),
	}

	return d, nil
}
