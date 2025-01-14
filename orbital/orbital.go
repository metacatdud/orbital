package orbital

import (
	"crypto/tls"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"orbital/pkg/logger"
	"strconv"
)

//go:embed web/*
var staticDir embed.FS

type Config struct {
	ApiServer       *Server
	WsServer        *WsConn
	Ip              string
	Port            int
	RootStoragePath string
}

type Orbital struct {
	client      *http.Server
	apiServer   *Server
	wsServer    *WsConn
	ip          string
	port        int
	rootStorage string
	log         *logger.Logger
}

func (n *Orbital) Start() error {
	if err := n.init(); err != nil {
		return err
	}

	n.log.Info("Serving Orbital dashboard", "addr", fmt.Sprintf("%s:%d", n.ip, n.port))

	//if err := n.client.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
	//	return fmt.Errorf("failed to start HTTP server: %w", err)
	//}

	if err := n.client.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

func (n *Orbital) init() error {
	staticFiles, err := fs.Sub(staticDir, "web")
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(staticFiles)))
	mux.Handle("/rpc/", n.apiServer)
	mux.Handle("/ws", n.wsServer)

	handler := CORSMiddleware(mux)

	//tlsCfg, err := tlsConfig(n.rootStorage)
	//if err != nil {
	//	return err
	//}

	n.client = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", n.ip, strconv.Itoa(n.port)),
		Handler: handler,
		//TLSConfig: tlsCfg,
	}

	return nil
}

func New(cfg Config) *Orbital {
	lg := logger.New(logger.LevelDebug, logger.FormatString)

	return &Orbital{
		apiServer:   cfg.ApiServer,
		wsServer:    cfg.WsServer,
		ip:          cfg.Ip,
		rootStorage: cfg.RootStoragePath,
		log:         lg,
		port:        cfg.Port,
	}
}

func tlsConfig(dataPath string) (*tls.Config, error) {
	certFile := fmt.Sprintf("%s/orbital/certs/server.crt", dataPath)
	keyFile := fmt.Sprintf("%s/orbital/certs/server.key", dataPath)

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificates: %v", err)
	}

	// Configure TLS
	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	return tlsCfg, nil
}
