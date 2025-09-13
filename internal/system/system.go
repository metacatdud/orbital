package system

import (
	"context"
	"orbital/config"
	"orbital/orbital"
	"orbital/pkg/cryptographer"
	"orbital/pkg/logger"
)

const (
	Domain              = "system"
	ActionKeepAlivePing = "keepAlivePing"
	ActionKeepAlivePong = "keepAlivePong"
	ActionWelcome       = "welcome"
)

type Dependencies struct {
	Log *logger.Logger
	Ws  *orbital.WsConn
}

type System struct {
	log *logger.Logger
	ws  *orbital.WsConn
}

func NewService(deps Dependencies) *System {
	return &System{
		log: deps.Log,
		ws:  deps.Ws,
	}
}

func (s *System) ConnectionKeepAlive(ctx context.Context, req ConnectionKeepAliveReq) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	sk, err := cryptographer.NewPrivateKeyFromHex(cfg.SecretKey)
	if err != nil {
		return err
	}

	meta := cryptographer.Metadata{
		Domain: Domain,
		Action: ActionKeepAlivePong,
	}

	msg, _ := cryptographer.Encode(sk, meta, ConnectionKeepAliveRes{
		Code: orbital.OK,
	})

	s.log.Debug("keep alive pong", "connId", req.ConnID)

	if err = s.ws.SendTo(ctx, req.ConnID, *msg); err != nil {
		return err
	}

	return nil
}
