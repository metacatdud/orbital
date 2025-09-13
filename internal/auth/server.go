package auth

import (
	"errors"
	"net/http"
	"orbital/config"
	"orbital/orbital"
	"orbital/pkg/cryptographer"
)

type authServiceServer struct {
	server  orbital.HTTPService
	service AuthService
}

func RegisterAuthServiceServer(server orbital.HTTPService, _ orbital.WsService, service AuthService) {
	handler := &authServiceServer{
		server:  server,
		service: service,
	}

	// Register middleware if any.
	// [!] These will be attached to all routes
	server.Use(
		MessageDecode(),
		ValidateRole(),
	)

	// Register routes
	server.Register(orbital.Route{
		ServiceName: "AuthService",
		ActionName:  "Auth",
		Handler:     handler.handleAuthentication,
		Method:      http.MethodPost,
	})

	server.Register(orbital.Route{
		ServiceName: "AuthService",
		ActionName:  "Check",
		Handler:     handler.handleCheckKey,
		Method:      http.MethodPost,
	})
}

func (s *authServiceServer) handleAuthentication(w http.ResponseWriter, r *http.Request) {

	publicKey, ok := r.Context().Value(cryptographer.PublicKeyCtxKey).(string)
	if !ok {
		s.server.OnError(w, r, errors.New("cannot decode body"))
		return
	}

	res, err := s.service.Auth(r.Context(), AuthReq{
		PublicKey: publicKey,
	})
	if err != nil {
		s.server.OnError(w, r, err)
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		s.server.OnError(w, r, err)
		return
	}

	sk, err := cryptographer.NewPrivateKeyFromHex(cfg.SecretKey)
	if err != nil {
		s.server.OnError(w, r, err)
		return
	}

	orbitalMessage, _ := cryptographer.Encode(sk, cryptographer.Metadata{
		Domain:        Domain,
		Action:        ActionLogin,
		CorrelationID: publicKey,
	}, res)

	if err = orbital.Encode(w, r, http.StatusOK, orbitalMessage); err != nil {
		s.server.OnError(w, r, err)
		return
	}
}

func (s *authServiceServer) handleCheckKey(w http.ResponseWriter, r *http.Request) {
	publicKey, ok := r.Context().Value(cryptographer.PublicKeyCtxKey).(string)
	if !ok {
		s.server.OnError(w, r, errors.New("cannot decode body"))
		return
	}

	res, err := s.service.Check(r.Context(), CheckReq{})
	if err != nil {
		s.server.OnError(w, r, err)
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		s.server.OnError(w, r, err)
		return
	}

	sk, err := cryptographer.NewPrivateKeyFromHex(cfg.SecretKey)
	if err != nil {
		s.server.OnError(w, r, err)
		return
	}

	orbitalMessage, _ := cryptographer.Encode(sk, cryptographer.Metadata{
		Domain:        Domain,
		Action:        ActionCheck,
		CorrelationID: publicKey,
	}, res)

	if err = orbital.Encode(w, r, http.StatusOK, orbitalMessage); err != nil {
		s.server.OnError(w, r, err)
		return
	}
}
