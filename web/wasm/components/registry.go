package components

import (
	"fmt"
	"orbital/web/wasm/orbital"
	"strings"
)

// init add components to registry here for later use.
// Use this only if the component can be mounted in several places
// otherwise prefer the constructor: NewXxxxxx
// TODO: Monitor this feature to see it makes sense. It feels like a bit of an overhead
func init() {
	RegisterComponent(LoginComponentRegKey, func(di *orbital.Dependency) Component {
		return NewLoginComponent(di)
	})
	RegisterComponent(DashboardComponentRegKey, func(di *orbital.Dependency) Component { return NewDashboardComponent(di) })
}

type RegKey string

func (key RegKey) String() string {
	return string(key)
}

func (key RegKey) WithExtra(sep string, parts ...any) RegKey {
	joined := fmt.Sprint(parts...)
	result := fmt.Sprintf("%s%s%s", key, sep, joined)
	result = strings.ReplaceAll(result, " ", "")
	return RegKey(result)
}

var registry = map[RegKey]func(di *orbital.Dependency) Component{}

func RegisterComponent(key RegKey, ctor func(di *orbital.Dependency) Component) {
	registry[key] = ctor
}

func LookupComponent(key RegKey, di *orbital.Dependency) (Component, error) {
	ctor, ok := registry[key]
	if !ok {
		return nil, fmt.Errorf("component not found: %s", key)
	}

	return ctor(di), nil
}
