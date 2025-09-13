package cmd

import (
	"orbital/config"
	"orbital/domain"
	"orbital/internal/apps"
	"orbital/internal/auth"
	"orbital/internal/machine"
	"orbital/internal/system"
	"orbital/orbital"
	"orbital/pkg/db"
	"orbital/pkg/logger"
	"orbital/pkg/prompt"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	debug bool
)

func newStartCmd() *cobra.Command {

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "start the node",
		RunE: func(cmd *cobra.Command, args []string) error {

			cmdHeader("start")

			var logLvl = logger.LevelError
			if debug {
				logLvl = logger.LevelDebug
			}

			log := logger.New(logLvl, logger.FormatString)

			cfg, err := config.LoadConfig()
			if err != nil {
				prompt.Err(prompt.NewLine("cannot load config: %s"), err.Error())
				return err
			}

			dbConn, err := setupDB(cfg)
			if err != nil {
				prompt.Err(prompt.NewLine("cannot setup db: %s"), err.Error())
				return err
			}

			// Repositories
			appRepo := domain.NewAppRepository(dbConn)
			userRepo := domain.NewUserRepository(dbConn)

			// Add TCP Server here

			apiSrv := orbital.NewServer(log)
			wsSrv := orbital.NewWsConn(log)

			// Prepare services
			authSvc := auth.NewService(auth.Dependencies{
				Log:      log,
				UserRepo: &userRepo,
				Ws:       wsSrv,
			})

			appsSvc := apps.NewService(apps.Dependencies{
				Log:     log,
				AppRepo: &appRepo,
			})

			machineSvc := machine.NewService(machine.Dependencies{
				Log: log,
				Ws:  wsSrv,
			})

			systemSvc := system.NewService(system.Dependencies{
				Log: log,
				Ws:  wsSrv,
			})

			// Register all service to server
			auth.RegisterAuthServiceServer(apiSrv, wsSrv, authSvc)
			apps.RegisterAppsServiceServer(apiSrv, wsSrv, appsSvc)
			machine.RegisterMachineServiceServer(apiSrv, wsSrv, machineSvc)
			system.RegisterSystemServiceServer(apiSrv, wsSrv, systemSvc)

			// Boot Orbital
			orbitalCfg := orbital.Config{
				ApiServer: apiSrv,
				WsServer:  wsSrv,
				Addr:      cfg.Addr,
				Cfg:       cfg,
				Logger:    log,
			}

			orbitalNode, err := orbital.New(orbitalCfg)
			if err != nil {
				return err
			}
			
			if err = orbitalNode.Start(); err != nil {
				return err
			}
			return nil
		},
	}

	startCmd.Flags().BoolVarP(&debug, "debug", "", false, "Debug mode")

	return startCmd
}

func setupDB(cfg *config.Config) (*db.DB, error) {
	dbPath := filepath.Join(cfg.OrbitalRootDir(), "data")

	dbConn, err := db.NewDB(dbPath)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}
