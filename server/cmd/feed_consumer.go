package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	cacheRedis "github.com/vitamin-nn/otus_architect_social/server/internal/cache/redis"
	"github.com/vitamin-nn/otus_architect_social/server/internal/config"
	"github.com/vitamin-nn/otus_architect_social/server/internal/consumer/feed"
	"github.com/vitamin-nn/otus_architect_social/server/internal/db/replication"
	"github.com/vitamin-nn/otus_architect_social/server/internal/queue/rabbit"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository/mysql"
)

func feedConsumerCmd(cfg *config.SocialConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "feed-consumer",
		Short: "starts feed consumer",
		Run: func(cmd *cobra.Command, args []string) {
			log.WithFields(cfg.Fields()).Info("starting feed consumer")

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

			defer func() {
				err := dbPool.Close()
				if err != nil {
					log.Fatalf("mysql close connect error: %v", err)
				}
			}()

			profileRepo := mysql.NewProfileRepo(dbPool)
			redisPool := cacheRedis.NewFeedCache(newRedisPool(context.Background(), cfg.Redis.Addr), maxFeedLen)

			rmq := rabbit.NewConsume(cfg.Rabbit.GetAddr(), "feed_consumer", cfg.Rabbit.ExchangeName, cfg.Rabbit.QueueName)

			log.Info("starting consuming feed message")
			err = rmq.Handle(feed.GetConsumerFunc(profileRepo, redisPool), 1)
			if err != nil {
				log.Fatalf("consumer handle error: %v", err)
			}

			interruptCh := make(chan os.Signal, 1)
			signal.Notify(interruptCh, os.Interrupt)
			signal.Notify(interruptCh, syscall.SIGTERM)
			log.Infof("graceful shutdown: %v", <-interruptCh)
			_, finish := context.WithTimeout(context.Background(), 5*time.Second)
			defer finish()
			rmq.Close()
		},
	}
}
