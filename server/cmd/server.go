package cmd

import (
	"context"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	jwtAuth "github.com/vitamin-nn/otus_architect_social/server/internal/auth/jwt"
	"github.com/vitamin-nn/otus_architect_social/server/internal/config"
	"github.com/vitamin-nn/otus_architect_social/server/internal/http"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository/mysql"
)

func serverCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Starts social network server",
		Run: func(cmd *cobra.Command, args []string) {
			log.WithFields(cfg.Fields()).Info("Starting antibruteforce service")

			//ctx, cancel := context.WithCancel(context.Background())
			jwt := jwtAuth.New(cfg.JWT.Secret, cfg.JWT.AccessLifeTime, cfg.JWT.RefreshLifeTime)
			dbConn, err := connDB(context.Background(), cfg.MySQL.GetDSN())
			if err != nil {
				log.Fatalf("mysql connect error: %v", err)
			}

			profileRepo := mysql.NewProfileRepo(dbConn)

			httpSrv := http.New(profileRepo, jwt, cfg.HTTPServer.WriteTimeout, cfg.HTTPServer.ReadTimeout)

			go func() {
				log.Info("Starting HTTP server")
				if err := httpSrv.Run(cfg.HTTPServer.Addr); err != nil {
					log.Fatal(err)
				}
			}()

			interruptCh := make(chan os.Signal, 1)
			signal.Notify(interruptCh, os.Interrupt)
			log.Infof("graceful shutdown: %v", <-interruptCh)
			ctx, finish := context.WithTimeout(context.Background(), 5*time.Second)
			//cancel()
			defer finish()
			if err := httpSrv.Shutdown(ctx); err != nil {
				log.Error("Error while shutdown")
			}

			err = dbConn.Close()
			if err != nil {
				log.Fatalf("mysql close connect error: %v", err)
			}
		},
	}
}
