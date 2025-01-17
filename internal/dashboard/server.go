package dashboard

import (
	"net/http"
	"orbital/orbital"
)

type dashboardServiceServer struct {
	httpServer orbital.HTTPService
	wsServer   orbital.WsService
	service    DashboardService
}

func RegisterDashboardServiceServer(httpServer orbital.HTTPService, wsServer orbital.WsService, service DashboardService) {
	handler := &dashboardServiceServer{
		httpServer: httpServer,
		wsServer:   wsServer,
		service:    service,
	}

	httpServer.Register(orbital.Route{
		ServiceName: "DashboardService",
		ActionName:  "RetrieveAllData",
		Handler:     handler.handleRetrieveAllData,
		Method:      http.MethodPut,
	})
}

func (s *dashboardServiceServer) handleRetrieveAllData(w http.ResponseWriter, r *http.Request) {

	var demo = map[string]string{
		"hello": "world from dashboard API",
	}

	if err := orbital.Encode(w, r, http.StatusOK, demo); err != nil {
		s.httpServer.OnError(w, r, err)
		return
	}
}
