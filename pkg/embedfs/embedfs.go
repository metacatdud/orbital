package embedfs

import (
	"embed"
	"fmt"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"net/http"
)

type embedFileSystem struct {
	http.Handler
	directory string
}

func (e embedFileSystem) Resolve(s string) string {
	fmt.Println("Resolve", s)
	return e.directory + s
}

func EmbedFileSystem(embedFS embed.FS, directory string) (app.ResourceResolver, error) {

	//dir, err := fs.Sub(embedFS, directory)
	//if err != nil {
	//	return nil, err
	//}

	httpFS := http.FileServer(http.FS(embedFS))
	return &embedFileSystem{
		Handler:   httpFS,
		directory: directory,
	}, nil
}
