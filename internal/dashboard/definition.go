package dashboard

import "context"

type DashboardService interface {
	RetrieveAllData(ctx context.Context, req RetrieveAllDataReq) (RetrieveAllDataResp, error)
}

type RetrieveAllDataReq struct {
}

type RetrieveAllDataResp struct{}
