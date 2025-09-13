package system

import (
	"context"
	"fmt"
	"log"
	"orbital/orbital"
)

type systemServiceServer struct {
	service SystemService
}

func RegisterSystemServiceServer(_ orbital.HTTPService, wsServer orbital.WsService, service SystemService) {
	h := &systemServiceServer{
		service: service,
	}

	wsServer.Register(orbital.Topic{
		Name:    "system/keepAlivePing",
		Handler: h.handleWsConnectionKeepAlive,
	})

	wsServer.Register(orbital.Topic{
		Name: "system/keepAlivePong",
		Handler: func(ctx context.Context, connID string, data []byte) {
			log.Printf("[!] keep alive pong: %s", connID)
		},
	})
}

func (h *systemServiceServer) handleWsConnectionKeepAlive(ctx context.Context, connID string, data []byte) {
	if err := h.service.ConnectionKeepAlive(ctx, ConnectionKeepAliveReq{
		ConnID: connID,
	}); err != nil {
		fmt.Printf("system: failed to keep alive: %v\n", err)
	}
}
