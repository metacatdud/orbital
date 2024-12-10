package orbital

import (
	"crypto/tls"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"orbital/pkg/logger"
)

//go:embed web/*
var staticDir embed.FS

type Config struct {
	ApiServer       *Server
	Ip              string
	RootStoragePath string
}

type Node struct {
	client      *http.Server
	apiServer   *Server
	ip          string
	rootStorage string
	log         *logger.Logger
}

func (n *Node) Start() error {
	if err := n.init(); err != nil {
		return err
	}

	n.log.Info("Serving Orbital dashboard", "addr", fmt.Sprintf("https://[%s]:8080", n.ip))

	if err := n.client.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

func (n *Node) init() error {
	staticFiles, err := fs.Sub(staticDir, "web")
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(staticFiles)))
	mux.Handle("/rpc/", n.apiServer)

	handler := CORSMiddleware(mux)

	tlsCfg, err := tlsConfig(n.rootStorage)
	if err != nil {
		return err
	}

	n.client = &http.Server{
		Addr:      fmt.Sprintf("[%s]:8080", n.ip),
		Handler:   handler,
		TLSConfig: tlsCfg,
	}

	return nil
}

func New(cfg Config) *Node {
	lg := logger.New(logger.LevelDebug, logger.FormatString)

	return &Node{
		apiServer:   cfg.ApiServer,
		ip:          cfg.Ip,
		rootStorage: cfg.RootStoragePath,
		log:         lg,
	}
}

func tlsConfig(dataPath string) (*tls.Config, error) {
	certFile := fmt.Sprintf("%s/server.crt", dataPath)
	keyFile := fmt.Sprintf("%s/server.key", dataPath)

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
