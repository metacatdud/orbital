package cmd

import (
	"encoding/hex"
	"orbital/pkg/cryptographer"
	"orbital/pkg/prompt"

	"github.com/spf13/cobra"
)

func newKeygenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keygen",
		Short: "Generate a new private key",
		RunE: func(cmd *cobra.Command, args []string) error {
			pk, sk, err := cryptographer.GenerateKeysPair()
			if err != nil {
				return err
			}

			prompt.OK("Secret key: %s\n", hex.EncodeToString(sk.Seed()))
			prompt.Info("Public key: %s\n", pk.ToHex())
			return nil
		},
	}

	return cmd
}
