package orbital

import (
	"orbital/orbital"
	"orbital/web/wasm/domain"
	"orbital/web/wasm/pkg/deps"
	"orbital/web/wasm/pkg/events"
	"orbital/web/wasm/pkg/transport"
)

type Machine struct {
	di     *deps.Dependency
	events *events.Event
	ws     *transport.WsConn
}

func (machine *Machine) init() {
	machine.ws.On("ws:orbital:machine", machine.wsMachines)
}

func (machine *Machine) wsMachines(data []byte) {
	machineRes, err := domain.NewMachineFromData(data)
	var machineResErr *orbital.ErrorResponse
	if err != nil {
		machineResErr = &orbital.ErrorResponse{
			Type: "auth.empty",
			Msg:  "private key cannot be empty",
		}

		machine.events.Emit("evt:machines:error", machineResErr)
		return
	}

	machine.events.Emit("evt:machines:update", machineRes)
}

func NewMachine(di *deps.Dependency) *Machine {
	machine := &Machine{
		di:     di,
		events: di.Events(),
		ws:     di.Ws(),
	}

	machine.init()

	return machine
}
