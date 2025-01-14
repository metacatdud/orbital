package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"orbital/orbital"
	"orbital/pkg/cryptographer"
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

	wsServer.Register(orbital.Topic{
		Name:    "dashboard.allData",
		Handler: handler.wsDummyHandler,
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

func (s *dashboardServiceServer) wsDummyHandler(connID string, data []byte) {
	var msg cryptographer.Message

	if err := json.Unmarshal(data, &msg); err != nil {
		fmt.Printf("Cannot marshal message: %v", msg)
		return
	}

	fmt.Printf("Got message from FE: %+v", string(data))
}
