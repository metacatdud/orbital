package orbital

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"orbital/config"
	"orbital/pkg/cryptographer"
	"orbital/pkg/logger"
)

//go:embed all:web/*
var staticDir embed.FS

type Config struct {
	ApiServer HTTPService
	WsServer  WsService
	Addr      string
	Cfg       *config.Config
	Logger    *logger.Logger
}

type Orbital struct {
	client    *http.Server
	apiServer HTTPService
	wsServer  WsService
	addr      string
	cfg       *config.Config
	log       *logger.Logger
}

func (n *Orbital) Start() error {

	staticFiles, err := fs.Sub(staticDir, "web")
	if err != nil {
		return err
	}

	httpFS := http.FileServer(http.FS(staticFiles))

	mux := http.NewServeMux()
	mux.Handle("/", fsFileHandlerMiddleware(httpFS))
	mux.Handle("/rpc/", n.apiServer)
	mux.Handle("/ws", n.wsServer)

	handler := corsMiddleware(mux)

	n.client = &http.Server{
		Addr:    n.addr,
		Handler: handler,
	}

	n.log.Info("Starting Orbital", "addr", n.addr)

	if err = n.client.ListenAndServe(); err != nil {
		return fmt.Errorf("%w:[%v]", ErrHttpListen, err)
	}

	return nil
}

func New(cfg Config) (*Orbital, error) {
	var lg *logger.Logger

	if cfg.Logger != nil {
		lg = cfg.Logger
	}

	if lg == nil {
		lg = logger.New(logger.LevelDebug, logger.FormatString)
	}

	sk, err := cryptographer.NewPrivateKeyFromHex(cfg.Cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	apiSrv := cfg.ApiServer
	apiSrv.SetSecretKey(sk)

	wsSrv := cfg.WsServer
	wsSrv.SetSecretKey(sk)

	return &Orbital{
		apiServer: apiSrv,
		wsServer:  wsSrv,
		addr:      cfg.Addr,
		cfg:       cfg.Cfg,
		log:       lg,
	}, nil
}
