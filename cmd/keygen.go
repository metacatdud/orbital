package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"orbital/pkg/cryptographer"
	"orbital/pkg/prompt"
)

var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate a new private key",
	RunE: func(cmd *cobra.Command, args []string) error {
		pk, sk, err := cryptographer.GenerateKeysPair()
		if err != nil {
			return err
		}

		cmdHeader("keygen")

		prompt.Err(prompt.NewLine("- Secret key: %s [DO NOT SHARE AND KEEP IT SAFE]"), sk.String())
		prompt.Info(prompt.NewLine("- Public key: %s"), pk.String())

		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(keygenCmd)
}
