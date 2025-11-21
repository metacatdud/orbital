package machine

import (
	"atomika.io/atomika/atomika"
)

type machineServiceServer struct {
	server  *atomika.HTTPService
	service MachineService
}

func RegisterMachineServiceServer(server *atomika.HTTPService, service MachineService) {
	_ = &machineServiceServer{
		server:  server,
		service: service,
	}
}
