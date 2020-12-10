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
	"github.com/vitamin-nn/otus_architect_social/server/internal/db/replication"
	"github.com/vitamin-nn/otus_architect_social/server/internal/http/server/profile"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository/mysql"
)

func serverCmd(cfg *config.SocialConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "profile",
		Short: "Starts social network server",
		Run: func(cmd *cobra.Command, args []string) {
			log.WithFields(cfg.Fields()).Info("Starting social server")

			jwt := jwtAuth.New(cfg.JWT.Secret, cfg.JWT.AccessLifeTime, cfg.JWT.RefreshLifeTime)

			dbConnMaster, err := connDB(context.Background(), cfg.MySQL.GetDSN())
			if err != nil {
				log.Fatalf("mysql master connect error: %v", err)
			}

			dbPool := replication.NewDBPool(dbConnMaster)

			for _, sl := range cfg.MySQL.SlavesDSN {
				dbConnSlave, err := connDB(context.Background(), sl)
				if err != nil {
					log.Fatalf("mysql slave connect error: %v", err)
				}
				dbPool.AddSlave(dbConnSlave)
			}

			profileRepo := mysql.NewProfileRepo(dbPool)

			httpSrv := profile.New(profileRepo, jwt, cfg.HTTPServer.WriteTimeout, cfg.HTTPServer.ReadTimeout)

			go func() {
				log.Info("Starting HTTP server")
				if err := httpSrv.Run(cfg.HTTPServer.GetAddr()); err != nil {
					log.Fatal(err)
				}
			}()

			interruptCh := make(chan os.Signal, 1)
			signal.Notify(interruptCh, os.Interrupt)
			log.Infof("graceful shutdown: %v", <-interruptCh)
			ctx, finish := context.WithTimeout(context.Background(), 5*time.Second)
			// cancel()
			defer finish()
			if err := httpSrv.Shutdown(ctx); err != nil {
				log.Error("Error while shutdown")
			}

			err = dbPool.Close()
			if err != nil {
				log.Fatalf("mysql close connect error: %v", err)
			}
		},
	}
}
