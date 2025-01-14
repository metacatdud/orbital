package auth

import (
	"net/http"
	"orbital/orbital"
)

type helloServiceServer struct {
	server  orbital.HTTPService
	service AuthService
}

func RegisterHelloServiceServer(server orbital.HTTPService, service AuthService) {
	handler := &helloServiceServer{server: server, service: service}

	server.Register(orbital.Route{
		ServiceName: "AuthService",
		ActionName:  "Auth",
		Handler:     handler.handleAuth,
		Method:      http.MethodPost,
	})
}

func (s *helloServiceServer) handleAuth(w http.ResponseWriter, r *http.Request) {

	var req AuthReq
	if err := orbital.Decode(r.Body, &req); err != nil {
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
