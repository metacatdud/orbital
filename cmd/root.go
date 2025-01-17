package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"orbital/pkg/prompt"
)

type Dependencies struct {
	FS fs.FS
}

var rootCmd = &cobra.Command{
	Use:   "orbital",
	Short: "orbital - container orchestration made simple",
	Long:  "orbital is a container orchestration system without the bloat",
}

func Execute(deps Dependencies) error {
	rootCmd.AddCommand(newInitCmd(deps))
	rootCmd.AddCommand(newKeygenCmd())
	rootCmd.AddCommand(newStartCmd())

	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}

func cmdHeader(section string) {
	prompt.Bold(prompt.ColorWhite, prompt.NewLine("Orbital |- %s"), section)
	prompt.Info(prompt.NewLine("----------------------------------------------------------------"))
	fmt.Println()
}
