package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"orbital/pkg/prompt"
)

var rootCmd = &cobra.Command{
	Use:   "orbital",
	Short: "orbital - container orchestration made simple",
	Long:  "orbital is a container orchestration system without the bloat",
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}

func cmdHeader(section string) {
	prompt.Bold(prompt.ColorWhite, prompt.NewLine("Orbital - %s"), section)
	prompt.Info(prompt.NewLine("----------------------------------------------------------------"))
	fmt.Println()
}
