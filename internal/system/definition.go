package system

import (
	"context"
	"orbital/orbital"
)

type SystemService interface {
	ConnectionKeepAlive(ctx context.Context, req ConnectionKeepAliveReq) error
}

type ConnectionKeepAliveReq struct {
	ConnID string `json:"connId"`
}

type ConnectionKeepAliveRes struct {
	Code  orbital.Code           `json:"code"`
	Error *orbital.ErrorResponse `json:"error,omitempty"`
}
