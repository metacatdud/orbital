package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"orbital/config"
	"orbital/domain"
	"orbital/internal/apps"
	"orbital/internal/auth"
	"orbital/internal/machine"
	"orbital/orbital"
	"orbital/pkg/db"
	"orbital/pkg/logger"
	"orbital/pkg/prompt"
	"path/filepath"
)

var (
	port int
)

func newStartCmd() *cobra.Command {

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "start the node",
		RunE: func(cmd *cobra.Command, args []string) error {

			cmdHeader("start")

			cfg, err := config.LoadConfig()
			if err != nil {
				return err
			}

			apiSrv, wsSrv, err := setupAPIServer(cfg)
			if err != nil {
				prompt.Err(prompt.NewLine("cannot start server: %s"), err.Error())
				return err
			}

			orbitalCfg := orbital.Config{
				ApiServer: apiSrv,
				WsServer:  wsSrv,
				Ip:        fmt.Sprintf("[%s]", cfg.BindIP),
				Cfg:       cfg,
				Port:      port,
			}

			orbitalNode := orbital.New(orbitalCfg)
			if err = orbitalNode.Start(); err != nil {
				return err
			}
			return nil
		},
	}

	startCmd.Flags().IntVarP(&port, "port", "p", 8080, "Service port")

	return startCmd
}

func setupAPIServer(cfg *config.Config) (*orbital.Server, *orbital.WsConn, error) {
	dbPath := filepath.Join(cfg.OrbitalRootDir(), "data")

	dbClient, err := db.NewDB(dbPath)
	if err != nil {
		return nil, nil, err
	}

	// Dependencies and repositories
	log := logger.New(logger.LevelDebug, logger.FormatString)

	userRepo := domain.NewUserRepository(dbClient)

	// Prepare servers
	apiSrv := orbital.NewServer(log)
	wsSrv := orbital.NewWsConn(log)

	// Prepare services
	authSvc := auth.NewService(auth.Dependencies{
		Log:      log,
		UserRepo: userRepo,
		Ws:       wsSrv,
	})

	appsSvc := apps.NewService(apps.Dependencies{
		Log: log,
	})

	machineSvc := machine.NewService(machine.Dependencies{
		Log: log,
		Ws:  wsSrv,
	})

	// Register all service to server
	auth.RegisterAuthServiceServer(apiSrv, wsSrv, authSvc)
	apps.RegisterAppsServiceServer(apiSrv, wsSrv, appsSvc)
	machine.RegisterMachineServiceServer(apiSrv, wsSrv, machineSvc)

	return apiSrv, wsSrv, nil
}
