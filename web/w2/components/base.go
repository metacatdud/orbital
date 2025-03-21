package components

import (
	"orbital/web/w2/pkg/dom"
	"syscall/js"
)

type BaseComp struct {
	docks map[string]js.Value
}

func (comp *BaseComp) GetContainer(name string) js.Value {
	container, ok := comp.docks[name]
	if !ok {
		return js.Null()
	}

	if container.IsNull() {
		return js.Null()
	}

	return container
}

func (comp *BaseComp) SetContainers(element js.Value) {
	if element.IsNull() {
		dom.ConsoleError("element is missing")
		return
	}

	dockingAreas := dom.QuerySelectorAllFromElement(element, `[data-dock]`)
	if len(dockingAreas) == 0 {
		return
	}

	for _, area := range dockingAreas {
		areaName := area.Get("dataset").Get("dock").String()
		comp.docks[areaName] = area
	}
}
