package apps

import (
	"encoding/json"
	"errors"
	"net/http"
	"orbital/config"
	"orbital/orbital"
	"orbital/pkg/cryptographer"
)

type appsServiceServer struct {
	server  orbital.HTTPService
	service AppsService
}

func RegisterAppsServiceServer(server orbital.HTTPService, _ orbital.WsService, service AppsService) {
	handler := &appsServiceServer{
		server:  server,
		service: service,
	}

	server.Register(orbital.Route{
		ServiceName: "AppsService",
		ActionName:  "List",
		Handler:     handler.handleList,
		Method:      http.MethodPost,
	})
}

func (s *appsServiceServer) handleList(w http.ResponseWriter, r *http.Request) {
	body, ok := r.Context().Value(cryptographer.BodyCtxKey).([]byte)
	if !ok {
		s.server.OnError(w, r, errors.New("cannot decode body"))
		return
	}

	var req ListReq
	if err := json.Unmarshal(body, &req); err != nil {
		s.server.OnError(w, r, err)
		return
	}

	res, err := s.service.List(r.Context(), req)
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
		Domain: Domain,
		Action: ActionList,
		Tags:   nil,
	}, res)

	if err = orbital.Encode(w, r, http.StatusOK, orbitalMessage); err != nil {
		s.server.OnError(w, r, err)
		return
	}
}
