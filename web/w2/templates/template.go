package templates

import (
	"bytes"
	"embed"
	"fmt"
	"golang.org/x/net/html"
	"html/template"
	"io/fs"
	"orbital/web/wasm/pkg/dom"
	"path/filepath"
	"strings"
)

//go:embed *
var templateFS embed.FS

type Registry struct {
	templates map[string]*template.Template
}

func NewRegistry() (*Registry, error) {
	r := &Registry{
		templates: make(map[string]*template.Template),
	}

	if err := r.loadTemplates(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Registry) Get(tplKey string) (*template.Template, error) {
	t, ok := r.templates[tplKey]
	if !ok {
		return nil, fmt.Errorf("template %s not found", tplKey)
	}

	return t, nil
}

func (r *Registry) loadTemplates() error {

	err := fs.WalkDir(templateFS, ".", func(path string, d fs.DirEntry, err error) error {

		if err != nil || d.IsDir() || filepath.Ext(path) != ".html" {
			return err
		}

		var content []byte
		content, err = templateFS.ReadFile(path)
		if err != nil {
			return err
		}

		baseKey := strings.TrimSuffix(strings.TrimPrefix(path, "./"), ".html")
		doc, err := html.Parse(bytes.NewReader(content))
		if err != nil {
			return err
		}

		var walk func(node *html.Node)
		walk = func(node *html.Node) {
			// Iterate and remove unneeded node and do proper trimming
			// to avoid unexpected behavior
			for c := node.FirstChild; c != nil; {
				next := c.NextSibling
				if c.Type == html.TextNode {
					trimmed := strings.TrimSpace(c.Data)
					if trimmed == "" {
						node.RemoveChild(c)
					} else {
						c.Data = trimmed
					}
				}

				if c.Type == html.ElementNode {
					walk(c)

					if c.FirstChild == nil && len(c.Attr) == 0 {
						node.RemoveChild(c)
					}
				}

				c = next
			}

			if node.Type == html.ElementNode && node.Data == "template" {
				var dataName string
				for _, attr := range node.Attr {
					if attr.Key == "data-name" {
						dataName = attr.Val
						break
					}
				}

				if dataName == "" {
					return
				}

				var buf bytes.Buffer
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					html.Render(&buf, c)
				}

				key := fmt.Sprintf("%s/%s", baseKey, dataName)
				r.templates[key] = template.Must(template.New(dataName).Parse(buf.String()))

				dom.ConsoleLog("- load:", key)
			}
		}

		walk(doc)

		return nil
	})

	return err

}
