package machine

import (
	"context"
	"orbital/orbital"
)

type MachineService interface {
	JobAllData(ctx context.Context, req AllDataReq) error
}

type AllDataReq struct {
}

type SystemInfo map[string]interface{}

type AllDataResp struct {
	SystemInfo *SystemInfo       `json:"systemInfo,omitempty"`
	Code       orbital.Code      `json:"code"`
	Error      map[string]string `json:"error,omitempty"`
}
