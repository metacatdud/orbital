package system

import (
	"context"
	"orbital/pkg/transport"
)

type SystemService interface {
	ConnectionKeepAlive(ctx context.Context, req ConnectionKeepAliveReq) error
}

type ConnectionKeepAliveReq struct {
	ConnID string `json:"connId"`
}

type ConnectionKeepAliveRes struct {
	Code  transport.Code         `json:"code"`
	Error *transport.ErrorResponse `json:"error,omitempty"`
}
