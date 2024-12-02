package cmd

import (
	"github.com/spf13/cobra"
	"orbital/internal/auth"
	"orbital/orbital"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start the node",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmdHeader("start")

		apiSrv := setupAPIServer()

		orbitalNode := orbital.New(apiSrv)
		if err := orbitalNode.Start(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func setupAPIServer() *orbital.Server {
	apiSrv := orbital.NewServer()

	helloService := auth.NewService()

	// Register all services to apiServer
	auth.RegisterHelloServiceServer(apiSrv, helloService)

	return apiSrv
}
