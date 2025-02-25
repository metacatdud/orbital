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

// SystemInfo basic aggregation struct
// TODO: Detail this for better experience latter on
type SystemInfo map[string]interface{}

type AllDataResp struct {
	SystemInfo *SystemInfo            `json:"systemInfo,omitempty"`
	Code       orbital.Code           `json:"code"`
	Error      *orbital.ErrorResponse `json:"error,omitempty"`
}
