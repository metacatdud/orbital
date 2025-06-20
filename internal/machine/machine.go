package machine

import (
	"context"
	"orbital/config"
	"orbital/orbital"
	"orbital/pkg/cryptographer"
	"orbital/pkg/jobber"
	"orbital/pkg/logger"
)

type Dependencies struct {
	Log *logger.Logger
	Ws  *orbital.WsConn
}

type Machine struct {
	jr  *jobber.Runner
	log *logger.Logger
	ws  *orbital.WsConn
}

func NewService(deps Dependencies) *Machine {
	m := &Machine{
		jr:  jobber.New(5),
		log: deps.Log,
		ws:  deps.Ws,
	}

	//m.jr.AddJob(10*time.Second, jobber.MaxRunInfinite, func() {
	//	_ = m.JobAllData(context.Background(), AllDataReq{})
	//})

	return m
}

func (service *Machine) JobAllData(ctx context.Context, req AllDataReq) error {

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	sk, err := cryptographer.NewPrivateKeyFromString(cfg.SecretKey)
	if err != nil {
		return err
	}

	meta := &cryptographer.Metadata{
		Domain: "machine",
		Action: "jobAllData",
	}

	body := &AllDataResp{}

	info, err := getInfo()
	if err != nil {
		body.Code = orbital.Internal
		body.Error = &orbital.ErrorResponse{
			Type: "machine.stats.err",
			Msg:  err.Error(),
		}

		msg, _ := cryptographer.Encode(sk, meta, body)
		service.ws.Broadcast(*msg)
	}

	cpu, err := getCPUInfo()
	if err != nil {
		body.Code = orbital.Internal
		body.Error = &orbital.ErrorResponse{
			Type: "machine.stats.err",
			Msg:  err.Error(),
		}

		msg, _ := cryptographer.Encode(sk, meta, body)
		service.ws.Broadcast(*msg)
	}

	mem, err := getMemInfo()
	if err != nil {
		body.Code = orbital.Internal
		body.Error = &orbital.ErrorResponse{
			Type: "machine.stats.err",
			Msg:  err.Error(),
		}
		msg, _ := cryptographer.Encode(sk, meta, body)
		service.ws.Broadcast(*msg)
	}

	netwk, err := getNetworkInfo()
	if err != nil {
		body.Code = orbital.Internal
		body.Error = &orbital.ErrorResponse{
			Type: "machine.stats.err",
			Msg:  err.Error(),
		}

		msg, _ := cryptographer.Encode(sk, meta, body)
		service.ws.Broadcast(*msg)
	}

	disk, err := getDiskInfo()
	if err != nil {
		body.Code = orbital.Internal
		body.Error = &orbital.ErrorResponse{
			Type: "machine.stats.err",
			Msg:  err.Error(),
		}
		msg, _ := cryptographer.Encode(sk, meta, body)
		service.ws.Broadcast(*msg)
	}

	stats := &SystemInfo{
		"info":  info,
		"cpu":   cpu,
		"disks": disk,
		"mem":   mem,
		"net":   netwk,
	}

	body.Code = orbital.OK
	body.SystemInfo = stats

	msg, _ := cryptographer.Encode(sk, meta, body)
	service.ws.Broadcast(*msg)

	return nil
}
