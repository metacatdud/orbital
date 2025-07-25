package templates

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"orbital/web/wasm/pkg/dom"
	"path/filepath"
	"regexp"
	"strings"
)

//go:embed *
var templateFS embed.FS

// templateReg regexp for collecting only `<template data-name="xyz"> ... </template>`
var templateReg = regexp.MustCompile(`(?s)<template\s+data-name="([^"]+)"[^>]*>(.*?)</template>`)

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
		matches := templateReg.FindAllSubmatch(content, -1)
		for _, match := range matches {
			tplKey := string(match[1])
			tplBody := string(match[2])
			key := fmt.Sprintf("%s/%s", baseKey, tplKey)

			var tpl *template.Template
			tpl, err = template.New(tplKey).Funcs(Funcs).Parse(tplBody)
			if err != nil {
				dom.ConsoleError("template key:", key, "err: ", err.Error())
				continue
			}

			r.templates[key] = tpl
			dom.ConsoleLog("- load:", key)
		}
		return nil
	})

	return err

}

// Funcs Helper Functions for HTML Templates
var Funcs = template.FuncMap{
	"contains": func(arr []string, v string) bool {
		for _, a := range arr {
			if a == v {
				return true
			}
		}
		return false
	},
}
