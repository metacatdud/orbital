package machine

import (
	"context"
	"orbital/pkg/transport"
)

type MachineService interface {
	JobAllData(ctx context.Context, req AllDataReq) error
}

type AllDataReq struct {
}

// SystemInfo basic aggregation struct
// TODO: Detail this with properties for better experience latter on
type SystemInfo map[string]any

type AllDataResp struct {
	SystemInfo *SystemInfo            `json:"systemInfo,omitempty"`
	Code       transport.Code         `json:"code"`
	Error      *transport.ErrorResponse `json:"error,omitempty"`
}
