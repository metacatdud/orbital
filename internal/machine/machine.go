package machine

import (
	"context"
	"orbital/config"
	"orbital/pkg/cryptographer"
	"orbital/pkg/jobber"
	"orbital/pkg/logger"
	"orbital/pkg/transport"

	"atomika.io/atomika/atomika"
)

type Dependencies struct {
	Log *logger.Logger
	Ws  atomika.WSDispatcher
}

type Machine struct {
	jr  *jobber.Runner
	log *logger.Logger
	ws  atomika.WSDispatcher
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

	sk, err := cryptographer.NewPrivateKeyFromHex(cfg.SecretKey)
	if err != nil {
		return err
	}

	meta := cryptographer.Metadata{
		Domain: "machine",
		Action: "jobAllData",
	}

	body := &AllDataResp{}

		info, err := getInfo()
		if err != nil {
			body.Code = transport.Internal
			body.Error = &transport.ErrorResponse{
				Type: "machine.stats.err",
				Msg:  err.Error(),
			}

			msg, _ := cryptographer.Encode(sk, meta, body)
			serialized, _ := msg.Serialize()
			service.ws.Broadcast(ctx, serialized)
		}

		cpu, err := getCPUInfo()
		if err != nil {
			body.Code = transport.Internal
			body.Error = &transport.ErrorResponse{
				Type: "machine.stats.err",
				Msg:  err.Error(),
			}

			msg, _ := cryptographer.Encode(sk, meta, body)
			serialized, _ := msg.Serialize()
			service.ws.Broadcast(ctx, serialized)
		}

		mem, err := getMemInfo()
		if err != nil {
			body.Code = transport.Internal
			body.Error = &transport.ErrorResponse{
				Type: "machine.stats.err",
				Msg:  err.Error(),
			}
			msg, _ := cryptographer.Encode(sk, meta, body)
			serialized, _ := msg.Serialize()
			service.ws.Broadcast(ctx, serialized)
		}

		netwk, err := getNetworkInfo()
		if err != nil {
			body.Code = transport.Internal
			body.Error = &transport.ErrorResponse{
				Type: "machine.stats.err",
				Msg:  err.Error(),
			}

			msg, _ := cryptographer.Encode(sk, meta, body)
			serialized, _ := msg.Serialize()
			service.ws.Broadcast(ctx, serialized)
		}

		disk, err := getDiskInfo()
		if err != nil {
			body.Code = transport.Internal
			body.Error = &transport.ErrorResponse{
				Type: "machine.stats.err",
				Msg:  err.Error(),
			}
			msg, _ := cryptographer.Encode(sk, meta, body)
			serialized, _ := msg.Serialize()
			service.ws.Broadcast(ctx, serialized)
		}

	stats := &SystemInfo{
		"info":  info,
		"cpu":   cpu,
		"disks": disk,
		"mem":   mem,
		"net":   netwk,
	}

	body.Code = transport.OK
	body.SystemInfo = stats

	msg, _ := cryptographer.Encode(sk, meta, body)
	serialized, _ := msg.Serialize()
	service.ws.Broadcast(ctx, serialized)

	return nil
}
