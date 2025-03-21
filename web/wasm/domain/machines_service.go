package domain

import (
	"orbital/orbital"
	orbital3 "orbital/web/wasm/orbital"
	orbital2 "orbital/web/wasm/pkg/events"
	"orbital/web/wasm/pkg/transport"
)

type MachineService struct {
	di     *orbital3.Dependency
	events *orbital2.Event
	ws     *transport.WsConn
}

func NewMachineService(di *orbital3.Dependency) *MachineService {
	machine := &MachineService{
		di:     di,
		events: di.Events(),
		ws:     di.Ws(),
	}

	machine.init()

	return machine
}

func (machine *MachineService) init() {
	machine.ws.On("ws:orbital:machine", machine.wsMachines)
}

func (machine *MachineService) wsMachines(data []byte) {
	machineRes, err := NewMachineFromData(data)
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
