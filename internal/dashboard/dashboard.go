package dashboard

import (
	"context"
)

type Dependencies struct {
}

type Dashboard struct {
}

func (d Dashboard) RetrieveAllData(ctx context.Context, req RetrieveAllDataReq) (RetrieveAllDataResp, error) {
	//TODO implement me
	panic("implement me")
}

func NewService(deps *Dependencies) *Dashboard {
	return &Dashboard{}
}
