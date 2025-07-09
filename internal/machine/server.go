package machine

import "orbital/orbital"

type machineServiceServer struct {
	server  orbital.HTTPService
	service MachineService
}

func RegisterMachineServiceServer(server orbital.HTTPService, wsServer orbital.WsService, service MachineService) {
	_ = &machineServiceServer{
		server:  server,
		service: service,
	}

	wsServer.Register(orbital.Topic{
		Name:    "",
		Handler: nil,
	})
}
