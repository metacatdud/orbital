package orbital

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"orbital/config"
	"orbital/pkg/logger"
	"strconv"
)

//go:embed all:web/*
var staticDir embed.FS

type Config struct {
	ApiServer *Server
	WsServer  *WsConn
	Ip        string
	Port      int
	Cfg       *config.Config
	Logger    *logger.Logger
}

type Orbital struct {
	client    *http.Server
	apiServer *Server
	wsServer  *WsConn
	ip        string
	port      int
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
		Addr:    fmt.Sprintf("%s:%s", n.ip, strconv.Itoa(n.port)),
		Handler: handler,
	}

	n.log.Info("Starting Orbital", "addr", fmt.Sprintf("%s:%d", n.ip, n.port))

	if err := n.client.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

func New(cfg Config) *Orbital {
	var lg *logger.Logger

	if cfg.Logger != nil {
		lg = cfg.Logger
	}

	if lg == nil {
		lg = logger.New(logger.LevelDebug, logger.FormatString)
	}

	return &Orbital{
		apiServer: cfg.ApiServer,
		wsServer:  cfg.WsServer,
		ip:        cfg.Ip,
		cfg:       cfg.Cfg,
		log:       lg,
		port:      cfg.Port,
	}
}
