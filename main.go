package main

import (
	"fmt"
	"orbital/cmd"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// TODO: Onboard: Admin Public key, storage path, http and server ports
// TODO: Instantiate an manager with HTTP
// TODO: Instantiate a SQLite handler
// TODO: Initiate a member list node
// TODO: Get machine's available resources
func run() error {
	return cmd.Execute()
}
