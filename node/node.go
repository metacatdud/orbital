package node

import (
	"embed"
	"errors"
	"fmt"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"net/http"
	"orbital/pkg/embedfs"
)

//go:embed web/*
var staticDir embed.FS

type Node struct {
	app *app.Handler
}

func (n *Node) Start() error {
	if err := n.init(); err != nil {
		return err
	}

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: n.app,
	}

	fmt.Println("Serving WASM dashboard at http://localhost:8080")
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

func (n *Node) init() error {
	//staticFiles, err := fs.Sub(staticDir, "web")
	//if err != nil {
	//	return err
	//}

	appDir, err := embedfs.EmbedFileSystem(staticDir, "web")
	if err != nil {
		return err
	}

	appHandler := &app.Handler{
		Name:      "Orbital",
		ShortName: "Orbital - docker orchestration",
		Resources: appDir,
	}

	n.app = appHandler

	//mux := http.NewServeMux()
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
	//	_, err = fs.Stat(staticFiles, r.URL.Path[1:])
	//	if err == nil {
	//		fileServer.ServeHTTP(w, r)
	//		return
	//	}
	//})
	//
	//n.httpServer = &http.Server{
	//	Addr:    ":8080",
	//	Handler: mux,
	//}

	return nil
}

func New() *Node {
	return &Node{}
}
