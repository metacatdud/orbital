package components

import "orbital/web/wasm/orbital"

const (
	TaskbarStartComponentRegKey RegKey = "taskbarStartComponent"
)

type TaskbarStartComp struct {
	*BaseComponent
}

func NewTaskbarStartComp(di *orbital.Dependency) *TaskbarStartComp {

	base := NewBaseComponent(di, TaskbarStartComponentRegKey, "orbital/taskbar/taskbarStartContent")
	comp := &TaskbarStartComp{
		BaseComponent: base,
	}

	return comp
}
