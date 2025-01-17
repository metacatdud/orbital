package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"orbital/config"
	"orbital/domain"
	"orbital/internal/auth"
	"orbital/internal/dashboard"
	"orbital/orbital"
	"orbital/pkg/db"
	"orbital/pkg/prompt"
	"path/filepath"
)

func newStartCmd() *cobra.Command {

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "start the node",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdHeader("start")

			cfg, err := config.LoadConfig("/etc/orbital/config.yaml")
			if err != nil {
				return err
			}

			apiSrv, wsSrv, err := setupAPIServer(cfg)
			if err != nil {
				prompt.Err(prompt.NewLine("cannot start server: %s"), err.Error())
				return err
			}

			orbitalCfg := orbital.Config{
				ApiServer:       apiSrv,
				WsServer:        wsSrv,
				Ip:              fmt.Sprintf("[%s]", cfg.BindIP),
				RootStoragePath: cfg.Datapath,
				Port:            8100,
			}

			orbitalNode := orbital.New(orbitalCfg)
			if err = orbitalNode.Start(); err != nil {
				return err
			}
			return nil
		},
	}

	return startCmd
}

func setupAPIServer(cfg *config.Config) (*orbital.Server, *orbital.WsConn, error) {
	dbPath := filepath.Join(cfg.OrbitalRootDir(), "data")

	dbClient, err := db.NewDB(dbPath)
	if err != nil {
		return nil, nil, err
	}

	// Dependencies
	userRepo := domain.NewUserRepository(dbClient)

	// Prepare server
	apiSrv := orbital.NewServer()
	wsSrv := orbital.NewWsConn()

	// Prepare services
	authService := auth.NewService(auth.Dependencies{
		UserRepo: userRepo,
		Ws:       wsSrv,
	})

	dashService := dashboard.NewService(nil)

	// Register all services to apiServer
	auth.RegisterHelloServiceServer(apiSrv, wsSrv, authService)
	dashboard.RegisterDashboardServiceServer(apiSrv, wsSrv, dashService)

	return apiSrv, wsSrv, nil
}
