package cmd

import (
	"github.com/spf13/cobra"
	"orbital/node"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start the node",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmdHeader("start")

		orbitalNode := node.New()
		if err := orbitalNode.Start(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
