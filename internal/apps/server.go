package apps

import (
	"encoding/json"
	"net/http"
	"orbital/config"
	"orbital/pkg/cryptographer"
	"orbital/pkg/transport"

	"atomika.io/atomika/atomika"
)

type appsServiceServer struct {
	server  *atomika.HTTPService
	service AppsService
}

func RegisterAppsServiceServer(server *atomika.HTTPService, service AppsService) {
	handler := &appsServiceServer{
		server:  server,
		service: service,
	}

	server.Register(atomika.Route{
		ServiceName: "AppsService",
		ActionName:  "List",
		Handler:     handler.handleList,
	})
}

func (s *appsServiceServer) handleList(w http.ResponseWriter, r *http.Request) {
	var req ListReq
	if err := transport.Decode(r.Body, &req); err != nil {
		_ = transport.Encode(w, r, http.StatusBadRequest, transport.ErrorResponse{Type: "apps.bad_request", Msg: err.Error()})
		return
	}

	res, err := s.service.List(r.Context(), req)
	if err != nil {
		_ = transport.Encode(w, r, http.StatusInternalServerError, transport.ErrorResponse{Type: "apps.error", Msg: err.Error()})
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		_ = transport.Encode(w, r, http.StatusInternalServerError, transport.ErrorResponse{Type: "apps.config", Msg: err.Error()})
		return
	}

	sk, err := cryptographer.NewPrivateKeyFromHex(cfg.SecretKey)
	if err != nil {
		_ = transport.Encode(w, r, http.StatusInternalServerError, transport.ErrorResponse{Type: "apps.key", Msg: err.Error()})
		return
	}

	orbitalMessage, _ := cryptographer.Encode(sk, cryptographer.Metadata{
		Domain: Domain,
		Action: ActionList,
		Tags:   nil,
	}, res)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	_ = enc.Encode(orbitalMessage)
}
