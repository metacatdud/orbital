package machine

import (
	"context"
	"orbital/orbital"
	"orbital/pkg/cryptographer"
	"orbital/pkg/jobber"
	"orbital/pkg/proto"
	"time"
)

type Dependencies struct {
	Ws *orbital.WsConn
}

type Machine struct {
	jr *jobber.Runner
	ws *orbital.WsConn
}

func (service *Machine) JobAllData(ctx context.Context, req AllDataReq) error {

	// DUMMY SERVER KEYS
	_, sk, _ := cryptographer.GenerateKeysPair()
	meta := &orbital.WsMetadata{
		Topic: "ws:orbital:machine",
	}

	body := &AllDataResp{}

	info, err := getInfo()
	if err != nil {
		body.Code = orbital.Internal
		body.Error = &orbital.ErrorResponse{
			Type: "machine.stats.err",
			Msg:  err.Error(),
		}

		msg, _ := proto.Encode(*sk, meta, body)
		service.ws.Broadcast(*msg)
	}

	cpu, err := getCPUInfo()
	if err != nil {
		body.Code = orbital.Internal
		body.Error = &orbital.ErrorResponse{
			Type: "machine.stats.err",
			Msg:  err.Error(),
		}

		msg, _ := proto.Encode(*sk, meta, body)
		service.ws.Broadcast(*msg)
	}

	mem, err := getMemInfo()
	if err != nil {
		body.Code = orbital.Internal
		body.Error = &orbital.ErrorResponse{
			Type: "machine.stats.err",
			Msg:  err.Error(),
		}
		msg, _ := proto.Encode(*sk, meta, body)
		service.ws.Broadcast(*msg)
	}

	netwk, err := getNetworkInfo()
	if err != nil {
		body.Code = orbital.Internal
		body.Error = &orbital.ErrorResponse{
			Type: "machine.stats.err",
			Msg:  err.Error(),
		}

		msg, _ := proto.Encode(*sk, meta, body)
		service.ws.Broadcast(*msg)
	}

	stats := &SystemInfo{
		"info": info,
		"cpu":  cpu,
		"mem":  mem,
		"net":  netwk,
	}

	body.Code = orbital.OK
	body.SystemInfo = stats

	msg, _ := proto.Encode(*sk, meta, body)
	service.ws.Broadcast(*msg)

	return nil
}

func NewService(deps Dependencies) *Machine {
	m := &Machine{
		jr: jobber.New(5),
		ws: deps.Ws,
	}

	m.jr.AddJob(30*time.Second, jobber.MaxRunInfinte, func() {
		_ = m.JobAllData(context.Background(), AllDataReq{})
	})

	return m
}
