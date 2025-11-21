package cmd

import (
	"net"
	"path/filepath"
	"time"

	"atomika.io/atomika/atomika"
	"orbital/config"
	"orbital/domain"
	"orbital/internal/apps"
	"orbital/internal/auth"
	"orbital/internal/machine"
	"orbital/internal/system"
	"orbital/pkg/db"
	"orbital/pkg/logger"
	"orbital/pkg/prompt"

	"github.com/spf13/cobra"
)

var debug bool

func newStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start the node",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdHeader("start")

			logLvl := logger.LevelError
			if debug {
				logLvl = logger.LevelDebug
			}
			log := logger.New(logLvl, logger.FormatString)

			cfg, err := config.LoadConfig()
			if err != nil {
				prompt.Err("\ncannot load config: %s", err.Error())
				return err
			}

			dbConn, err := setupDB(cfg)
			if err != nil {
				prompt.Err("\ncannot setup db: %s", err.Error())
				return err
			}

			appRepo := domain.NewAppRepository(dbConn)
			userRepo := domain.NewUserRepository(dbConn)

			httpSvc, err := buildHTTPService(cfg)
			if err != nil {
				prompt.Err("\ncannot setup http: %s", err.Error())
				return err
			}

			wsDisp := httpSvc.Dispatcher()

			authSvc := auth.NewService(auth.Dependencies{Log: log, UserRepo: &userRepo, Ws: wsDisp})
			appsSvc := apps.NewService(apps.Dependencies{Log: log, AppRepo: &appRepo})
			machineSvc := machine.NewService(machine.Dependencies{Log: log, Ws: wsDisp})
			systemSvc := system.NewService(system.Dependencies{Log: log, Ws: wsDisp})

			auth.RegisterAuthServiceServer(httpSvc, authSvc)
			apps.RegisterAppsServiceServer(httpSvc, appsSvc)
			machine.RegisterMachineServiceServer(httpSvc, machineSvc)
			system.RegisterSystemServiceServer(httpSvc, httpSvc, systemSvc)

			runtime := atomika.New()
			runtime.RegisterServices([]atomika.Service{httpSvc})

			return runtime.Boot()
		},
	}

	cmd.Flags().BoolVarP(&debug, "debug", "", false, "Debug mode")
	return cmd
}

func setupDB(cfg *config.Config) (*db.DB, error) {
	dbPath := filepath.Join(cfg.OrbitalRootDir(), "data")
	return db.NewDB(dbPath)
}

func buildHTTPService(cfg *config.Config) (*atomika.HTTPService, error) {
	port := extractPort(cfg.Addr)
	httpCfg := &atomika.CfgHttp{
		Port:     port,
		BasePath: "/rpc/",
		WWW:      "orbital/web",
		Websocket: &atomika.CfgWebsocket{
			Enable: true,
			Path:   "/ws",
		},
	}

	svc, err := atomika.NewHTTPService(httpCfg)
	if err != nil {
		return nil, err
	}

	svc.ConfigureWebsocket(httpCfg.Websocket, atomika.WithPingPong(), atomika.WithKeepAlive(60*time.Second))
	return svc, nil
}

func extractPort(addr string) string {
	_, port, err := net.SplitHostPort(addr)
	if err != nil || port == "" {
		return "8080"
	}
	return port
}
