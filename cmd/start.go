package cmd

import (
	"github.com/spf13/cobra"
	"orbital/config"
	"orbital/domain"
	"orbital/internal/auth"
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

			apiSrv, err := setupAPIServer(cfg)
			if err != nil {
				prompt.Err(prompt.NewLine("cannot start server: %s"), err.Error())
				return err
			}

			orbitalCfg := orbital.Config{
				ApiServer:       apiSrv,
				Ip:              cfg.BindIP,
				RootStoragePath: cfg.Datapath,
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

func setupAPIServer(cfg *config.Config) (*orbital.Server, error) {
	dbPath := filepath.Join(cfg.OrbitalRootDir(), "data")

	dbClient, err := db.NewDB(dbPath)
	if err != nil {
		return nil, err
	}

	// Dependencies
	userRepo := domain.NewUserRepository(dbClient)

	// Prepare server
	apiSrv := orbital.NewServer()

	// Prepare services
	authService := auth.NewService(auth.Dependencies{
		UserRepo: userRepo,
	})

	// Register all services to apiServer
	auth.RegisterHelloServiceServer(apiSrv, authService)

	return apiSrv, nil
}
