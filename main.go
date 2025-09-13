package main

import (
	"embed"
	"fmt"
	"orbital/cmd"
	"orbital/pkg/cryptographer"
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
	sk, err := cryptographer.NewPrivateKeyFromHex("4819b1b3acf279cb138da57e5debc274f1fa11cc9d17b35d876d6447baa6e1b4")
	if err != nil {
		return fmt.Errorf("cannot create private key: %w", err)
	}

	body := []byte("Hello World")
	msg := &cryptographer.Message{
		V:         0,
		Timestamp: cryptographer.Now(),
		Metadata:  cryptographer.Metadata{},
		Body:      body,
	}

	if err = msg.Sign(sk.Seed()); err != nil {
		return fmt.Errorf("cannot sign message: %w", err)
	}

	return cmd.Execute(cmd.Dependencies{
		FS: resourcesDir,
	})
}
