package auth

import (
	"encoding/json"
	"net/http"
	"orbital/config"
	"orbital/pkg/cryptographer"
	"orbital/pkg/transport"

	"atomika.io/atomika/atomika"
)

type authServiceServer struct {
	server  *atomika.HTTPService
	service AuthService
}

func RegisterAuthServiceServer(server *atomika.HTTPService, service AuthService) {
	handler := &authServiceServer{
		server:  server,
		service: service,
	}

	server.Register(atomika.Route{
		ServiceName: "AuthService",
		ActionName:  "Auth",
		Handler:     handler.handleAuthentication,
	})

	server.Register(atomika.Route{
		ServiceName: "AuthService",
		ActionName:  "Check",
		Handler:     handler.handleCheckKey,
	})
}

func (s *authServiceServer) handleAuthentication(w http.ResponseWriter, r *http.Request) {
	var req AuthReq
	if err := transport.Decode(r.Body, &req); err != nil {
		_ = transport.Encode(w, r, http.StatusBadRequest, transport.ErrorResponse{Type: "auth.bad_request", Msg: err.Error()})
		return
	}

	res, err := s.service.Auth(r.Context(), req)
	if err != nil {
		_ = transport.Encode(w, r, http.StatusInternalServerError, transport.ErrorResponse{Type: "auth.error", Msg: err.Error()})
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		_ = transport.Encode(w, r, http.StatusInternalServerError, transport.ErrorResponse{Type: "auth.config", Msg: err.Error()})
		return
	}

	sk, err := cryptographer.NewPrivateKeyFromHex(cfg.SecretKey)
	if err != nil {
		_ = transport.Encode(w, r, http.StatusInternalServerError, transport.ErrorResponse{Type: "auth.key", Msg: err.Error()})
		return
	}

	orbitalMessage, _ := cryptographer.Encode(sk, cryptographer.Metadata{
		Domain:        Domain,
		Action:        ActionLogin,
		CorrelationID: req.PublicKey,
	}, res)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	_ = enc.Encode(orbitalMessage)
}

func (s *authServiceServer) handleCheckKey(w http.ResponseWriter, r *http.Request) {
	var req CheckReq
	_ = transport.Decode(r.Body, &req)

	res, err := s.service.Check(r.Context(), req)
	if err != nil {
		_ = transport.Encode(w, r, http.StatusInternalServerError, transport.ErrorResponse{Type: "auth.error", Msg: err.Error()})
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		_ = transport.Encode(w, r, http.StatusInternalServerError, transport.ErrorResponse{Type: "auth.config", Msg: err.Error()})
		return
	}

	sk, err := cryptographer.NewPrivateKeyFromHex(cfg.SecretKey)
	if err != nil {
		_ = transport.Encode(w, r, http.StatusInternalServerError, transport.ErrorResponse{Type: "auth.key", Msg: err.Error()})
		return
	}

	orbitalMessage, _ := cryptographer.Encode(sk, cryptographer.Metadata{
		Domain:        Domain,
		Action:        ActionCheck,
		CorrelationID: "",
	}, res)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	_ = enc.Encode(orbitalMessage)
}
