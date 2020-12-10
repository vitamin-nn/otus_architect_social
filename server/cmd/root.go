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
