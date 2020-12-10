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
	"github.com/vitamin-nn/otus_architect_social/server/internal/db/sharding"
	"github.com/vitamin-nn/otus_architect_social/server/internal/http/server/messenger"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository/mysql"
)

func messengerCmd(cfg *config.MessengerConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "messenger",
		Short: "Starts messenger server",
		Run: func(cmd *cobra.Command, args []string) {
			log.WithFields(cfg.Fields()).Info("Starting messenger server")

			jwt := jwtAuth.New(cfg.JWT.Secret, cfg.JWT.AccessLifeTime, cfg.JWT.RefreshLifeTime)

			dbShardPool := sharding.NewDBShardPool()

			for _, sl := range cfg.MySQL.ShardsDSN {
				dbConnShard, err := connDB(context.Background(), sl)
				if err != nil {
					log.Fatalf("mysql slave connect error: %v", err)
				}
				dbShardPool.AddShard(dbConnShard)
			}

			messengerRepo := mysql.NewMessengerRepo(dbShardPool)

			httpSrv := messenger.NewMessenger(messengerRepo, jwt, cfg.HTTPServer.WriteTimeout, cfg.HTTPServer.ReadTimeout)

			go func() {
				log.Info("Starting HTTP messenger server")
				if err := httpSrv.Run(cfg.HTTPServer.GetAddr()); err != nil {
					log.Fatal(err)
				}
			}()

			interruptCh := make(chan os.Signal, 1)
			signal.Notify(interruptCh, os.Interrupt)
			log.Infof("graceful shutdown: %v", <-interruptCh)
			ctx, finish := context.WithTimeout(context.Background(), 5*time.Second)

			defer finish()
			if err := httpSrv.Shutdown(ctx); err != nil {
				log.Error("Error while shutdown")
			}

			err := dbShardPool.Close()
			if err != nil {
				log.Fatalf("mysql close connect error: %v", err)
			}
		},
	}
}
