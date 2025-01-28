package main

import (
	"embed"
	"fmt"
	"orbital/cmd"
	"os"
)

//go:embed resources/*
var resourcesDir embed.FS

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	return cmd.Execute(cmd.Dependencies{
		FS: resourcesDir,
	})
}
