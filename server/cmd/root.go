package cmd

import (
	"context"
	"database/sql"
	"log"

	"github.com/spf13/cobra"
	"github.com/vitamin-nn/otus_architect_social/server/internal/config"
	"github.com/vitamin-nn/otus_architect_social/server/internal/logger"
)

func Execute() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config file read error: %v", err)
	}

	err = logger.Init(cfg.Log)
	if err != nil {
		log.Fatalf("initialize logger error: %v", err)
	}

	rootCmd := &cobra.Command{
		Use:   "social",
		Short: "Social network",
	}

	rootCmd.AddCommand(serverCmd(cfg))

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
