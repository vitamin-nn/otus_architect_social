package cmd

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/spf13/cobra"
	"github.com/vitamin-nn/otus_architect_social/server/internal/config"
	"github.com/vitamin-nn/otus_architect_social/server/internal/db/replication"
	"github.com/vitamin-nn/otus_architect_social/server/internal/logger"
)

const maxFeedLen = 1000

func Execute() {
	cfgSocial, err := config.LoadSocial()
	if err != nil {
		log.Fatalf("config file read error: %v", err)
	}

	err = logger.Init(cfgSocial.Log)
	if err != nil {
		log.Fatalf("initialize logger error: %v", err)
	}

	rootCmd := &cobra.Command{
		Use:   "social",
		Short: "Social network",
	}
	rootCmd.AddCommand(serverCmd(cfgSocial))

	cfgMsg, err := config.LoadMessenger()
	if err != nil {
		log.Fatalf("config file read error: %v", err)
	}
	rootCmd.AddCommand(messengerCmd(cfgMsg))

	rootCmd.AddCommand(feedConsumerCmd(cfgSocial))

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("execute cmd: %v", err)
	}
}

func connDB(ctx context.Context, dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.Stats()

	return db, db.PingContext(ctx)
}

func NewDBPool(ctx context.Context, addrMaster string, addrSlaveList []string) *replication.DBReplPool {
	dbConnMaster, err := connDB(ctx, addrMaster)
	if err != nil {
		log.Fatalf("mysql master connect error: %v", err)
	}

	dbPool := replication.NewDBPool(dbConnMaster)

	for _, sl := range addrSlaveList {
		dbConnSlave, err := connDB(ctx, sl)
		if err != nil {
			log.Fatalf("mysql slave connect error: %v", err)
		}
		dbPool.AddSlave(dbConnSlave)
	}

	return dbPool
}

func newRedisPool(ctx context.Context, addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		DialContext: func(_ context.Context) (redis.Conn, error) {
			return redis.DialContext(ctx, "tcp", addr)
		},
	}
}
