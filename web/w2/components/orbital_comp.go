package components

import (
	"bytes"
	"errors"
	"orbital/web/w2/orbital"
	"orbital/web/w2/pkg/dom"
	"syscall/js"
)

const (
	OrbitalCompKey = "orbitalComponent"
)

type OrbitalComp struct {
	BaseComp
	di *orbital.Dependency

	element js.Value

	// Resident components
	taskbarComp orbital.Component
	overlayComp orbital.Component
	desktopComp orbital.Component
}

// Implementation checklist
var _ orbital.ContainerComponent = (*OrbitalComp)(nil)

func NewOrbitalComp(di *orbital.Dependency) *OrbitalComp {
	comp := &OrbitalComp{
		BaseComp: BaseComp{docks: make(map[string]js.Value)},
		di:       di,
	}

	return comp
}

func (comp *OrbitalComp) ID() string {
	return OrbitalCompKey
}

func (comp *OrbitalComp) Namespace() string {
	return "orbital/main/orbital"
}

func (comp *OrbitalComp) Mount(container *js.Value) error {
	if !container.Truthy() {
		return errors.New("container does not exist")
	}

	loadingElem := dom.QuerySelector("#loading")
	if !loadingElem.IsNull() {
		dom.RemoveElement(loadingElem)
	}

	if comp.element.IsNull() {
		return errors.New("element is missing")
	}

	dom.AppendChild(*container, comp.element)

	comp.mountTaskbar()
	comp.mountOverlay()

	return nil
}

func (comp *OrbitalComp) Unmount() error {
	return nil
}

func (comp *OrbitalComp) Render() error {
	tpl, err := comp.di.TplRegistry().Get(comp.Namespace())
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, nil); err != nil {
		return err
	}

	comp.element = dom.CreateElementFromString(buf.String())

	comp.SetContainers(comp.element)

	return nil
}

func (comp *OrbitalComp) mountOverlay() {
	container := comp.GetContainer("overlay")
	if container.IsNull() {
		dom.ConsoleError("overlay component container is null", comp.ID())
		return
	}
	dom.SetInnerHTML(container, "")

	var err error
	comp.overlayComp, err = orbital.Lookup[*OverlayComp](OverlayCompKey)
	if err != nil {
		dom.ConsoleError("overlay component cannot be created", comp.ID())
		return
	}

	if err = comp.overlayComp.Render(); err != nil {
		dom.ConsoleError("overlay component cannot be rendered", err.Error(), comp.ID())
		return
	}

	if err = comp.overlayComp.Mount(&container); err != nil {
		dom.ConsoleError("overlay component cannot be mounted", err.Error(), comp.ID())
		return
	}
}

func (comp *OrbitalComp) mountTaskbar() {
	container := comp.GetContainer("taskbar")
	if container.IsNull() {
		dom.ConsoleError("taskbar component container is null", comp.ID())
		return
	}
	dom.SetInnerHTML(container, "")

	var err error
	comp.taskbarComp, err = orbital.Lookup[*TaskbarComp](TaskbarCompKey)
	if err != nil {
		dom.ConsoleError("taskbar component cannot be created", comp.ID())
		return
	}

	if err = comp.taskbarComp.Render(); err != nil {
		dom.ConsoleError("taskbar component cannot be rendered", err.Error(), comp.ID())
		return
	}

	if err = comp.taskbarComp.Mount(&container); err != nil {
		dom.ConsoleError("taskbar component cannot be mounted", err.Error(), comp.ID())
		return
	}
}
