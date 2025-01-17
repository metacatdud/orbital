package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"orbital/orbital"
	"orbital/pkg/proto"
)

type authServiceServer struct {
	server  orbital.HTTPService
	service AuthService
}

func RegisterHelloServiceServer(server orbital.HTTPService, wsServer orbital.WsService, service AuthService) {
	handler := &authServiceServer{server: server, service: service}

	server.Register(orbital.Route{
		ServiceName: "AuthService",
		ActionName:  "Auth",
		Handler:     handler.handleAuth,
		Method:      http.MethodPost,
	})

	wsServer.Register(orbital.Topic{
		Name:    "orbital.authentication",
		Handler: handler.wsOrbitalAuthentication,
	})
}

func (s *authServiceServer) handleAuth(w http.ResponseWriter, r *http.Request) {

	var protoMessage proto.Message
	if err := orbital.Decode(r.Body, &protoMessage); err != nil {
		s.server.OnError(w, r, err)
		return
	}

	var req AuthReq
	if err := proto.Decode(protoMessage, &req, nil); err != nil {
		s.server.OnError(w, r, err)
		return
	}

	res, err := s.service.Auth(r.Context(), req)
	if err != nil {
		s.server.OnError(w, r, err)
		return
	}

	if err = orbital.Encode(w, r, http.StatusOK, res); err != nil {
		s.server.OnError(w, r, err)
		return
	}
}

func (s *authServiceServer) wsOrbitalAuthentication(connID string, data []byte) {
	var protoMessage proto.Message

	if err := json.Unmarshal(data, &protoMessage); err != nil {
		fmt.Printf("(ConnID: %s) Cannot marshal message", connID)
		return
	}

	var req WsAuthReq
	if err := proto.Decode(protoMessage, &req, nil); err != nil {
		fmt.Printf("(ConnID: %s) Cannot unmarshal message", connID)
		return
	}

	if err := s.service.WsAuth(context.Background(), connID, req); err != nil {
		fmt.Printf("(ConnID: %s) Cannot handle auth request", connID)
		return
	}

	return
}
