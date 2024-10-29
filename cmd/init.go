package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"orbital/pkg/prompt"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize orbital node",
	RunE: func(cmd *cobra.Command, args []string) error {

		cmdHeader("init")

		secretKey, _ := cmd.Flags().GetString("sk")
		ip, _ := cmd.Flags().GetString("ip")
		dataPath, _ := cmd.Flags().GetString("datapath")

		if secretKey == "" || ip == "" || dataPath == "" {
			return errors.New("secret key, ip and datapath cannot be empty")
		}

		prompt.Info(prompt.NewLine("- Secret key: %s"), secretKey)
		prompt.Info(prompt.NewLine("- IP: %s"), ip)
		prompt.Info(prompt.NewLine("- Data storage path: %s"), dataPath)

		fmt.Println()

		return nil
	},
}

func init() {
	initCmd.Flags().String("sk", "", "Secret key for node communication. Use keygen command to generate")
	initCmd.Flags().String("ip", "", "Node binding ip")
	initCmd.Flags().String("datapath", "", "Orbital data storage path")
	rootCmd.AddCommand(initCmd)
}
