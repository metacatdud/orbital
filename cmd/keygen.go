package cmd

import (
	"encoding/hex"
	"fmt"
	"orbital/pkg/cryptographer"
	"orbital/pkg/prompt"

	"github.com/spf13/cobra"
)

func newKeygenCmd() *cobra.Command {

	keygenCmd := &cobra.Command{
		Use:   "keygen",
		Short: "Generate a new private key",
		RunE: func(cmd *cobra.Command, args []string) error {
			pk, sk, err := cryptographer.GenerateKeysPair()
			if err != nil {
				return err
			}

			cmdHeader("keygen")

			prompt.Err(prompt.NewLine("- Secret key: %s [DO NOT SHARE AND KEEP IT SAFE]"), hex.EncodeToString(sk.Seed()))
			prompt.Info(prompt.NewLine("- Public key: %s"), pk.ToHex())

			fmt.Println()
			return nil
		},
	}

	return keygenCmd
}
