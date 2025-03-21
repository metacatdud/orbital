package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"orbital/orbital"
	"orbital/pkg/proto"
)

type authServiceServer struct {
	server  orbital.HTTPService
	service AuthService
}

func RegisterAuthServiceServer(server orbital.HTTPService, wsServer orbital.WsService, service AuthService) {
	handler := &authServiceServer{
		server:  server,
		service: service,
	}

	wsServer.Register(orbital.Topic{
		Name:    "orbital:authenticate",
		Handler: handler.wsOrbitalAuthentication,
	})
}

func (s *authServiceServer) wsOrbitalAuthentication(connID string, data []byte) {
	var protoMessage proto.Message
	if err := json.Unmarshal(data, &protoMessage); err != nil {
		fmt.Printf("(ConnID: %s) Cannot marshal message", connID)
		return
	}

	var req AuthReq
	if err := proto.Decode(protoMessage, &req, nil); err != nil {
		fmt.Printf("(ConnID: %s) Cannot unmarshal message", connID)
		return
	}

	if err := s.service.Auth(context.Background(), connID, req); err != nil {
		fmt.Printf("(ConnID: %s) Cannot handle auth request", connID)
		return
	}

	return
}
