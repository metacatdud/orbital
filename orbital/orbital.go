package orbital

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed web/*
var staticDir embed.FS

type Node struct {
	client    *http.Server
	apiServer *Server
}

func (n *Node) Start() error {
	if err := n.init(); err != nil {
		return err
	}

	fmt.Println("Serving WASM dashboard at http://localhost:8080")
	if err := n.client.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

	//mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//
	//	// Skip favicon requests
	//	if strings.Contains(r.URL.Path, "favicon.ico") {
	//		w.WriteHeader(http.StatusNoContent)
	//		return
	//	}
	//
	//	var data []byte
	//
	//	if r.URL.Path == "/" {
	//
	//		data, err = fs.ReadFile(staticFiles, "index.html")
	//		if err != nil {
	//			http.Error(w, "index.html not found", http.StatusNotFound)
	//			return
	//		}
	//		w.Header().Set("Content-Type", "text/html")
	//		w.Write(data)
	//		return
	//	}
	//
	//	// Otherwise, serve files from the embedded web folder
	//	http.FileServer(http.FS(staticFiles)).ServeHTTP(w, r)
	//})

	n.client = &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	return nil
}

func New(apiServer *Server) *Node {
	return &Node{
		apiServer: apiServer,
	}
}
