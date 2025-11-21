package cmd

import (
	"io/fs"

	"github.com/spf13/cobra"
	"orbital/pkg/prompt"
)

type Dependencies struct {
	FS fs.FS
}

func Execute(deps Dependencies) error {
	cmd := &cobra.Command{
		Use:           "orbital",
		Short:         "orbital CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(newInitCmd(deps))
	cmd.AddCommand(newUpdateCmd(deps))
	cmd.AddCommand(newKeygenCmd())
	cmd.AddCommand(newStartCmd())

	if _, err := cmd.ExecuteC(); err != nil {
		return err
	}
	return nil
}

func cmdHeader(section string) {
	prompt.Bold(prompt.ColorWhite, "\nOrbital |- %s\n", section)
	prompt.Info("----------------------------------------------------------------\n")
}
