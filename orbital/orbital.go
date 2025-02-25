package orbital

import (
	"crypto/tls"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"orbital/pkg/logger"
	"strconv"
	"strings"
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

	n.log.Info("Starting Orbital", "addr", fmt.Sprintf("%s:%d", n.ip, n.port))

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

	httpFS := http.FileServer(http.FS(staticFiles))

	mux := http.NewServeMux()
	mux.Handle("/", fsFileHandlerMiddleware(httpFS))
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

func fsFileHandlerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.HasSuffix(r.URL.Path, ".wasm") {
			w.Header().Set("Content-Type", "application/wasm")

			// This WASM should be considered dev build
			// and not cached
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		}

		if strings.HasSuffix(r.URL.Path, ".wasm.br") {
			w.Header().Set("Content-Type", "application/wasm")
			w.Header().Set("Content-Encoding", "br")

		}

		h.ServeHTTP(w, r)
	})
}
