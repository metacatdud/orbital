package system

import (
	"context"
	"fmt"
	"log"

	"atomika.io/atomika/atomika"
)

type systemServiceServer struct {
	service SystemService
}

func RegisterSystemServiceServer(_ *atomika.HTTPService, wsServer *atomika.HTTPService, service SystemService) {
	h := &systemServiceServer{
		service: service,
	}

	wsServer.RegisterTopic(atomika.WSTopic{
		Name: "system/keepAlivePing",
		Handler: h.handleWsConnectionKeepAlive,
	})

	wsServer.RegisterTopic(atomika.WSTopic{
		Name: "system/keepAlivePong",
		Handler: func(ctx context.Context, connID string, _ []byte) {
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
