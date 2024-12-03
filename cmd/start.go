package cmd

import (
	"github.com/spf13/cobra"
	"orbital/config"
	"orbital/internal/auth"
	"orbital/orbital"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start the node",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmdHeader("start")

		cfg, err := config.LoadConfig("/etc/orbital/config.yaml")
		if err != nil {
			return err
		}

		apiSrv := setupAPIServer()

		orbitalCfg := orbital.Config{
			ApiServer:       apiSrv,
			Ip:              cfg.BindIP,
			RootStoragePath: cfg.Datapath,
		}

		orbitalNode := orbital.New(orbitalCfg)
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
