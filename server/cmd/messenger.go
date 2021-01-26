package cmd

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vitamin-nn/otus_architect_social/server/internal/auth"
	jwtAuth "github.com/vitamin-nn/otus_architect_social/server/internal/auth/jwt"
	"github.com/vitamin-nn/otus_architect_social/server/internal/config"
	"github.com/vitamin-nn/otus_architect_social/server/internal/db/sharding"
	authMiddleware "github.com/vitamin-nn/otus_architect_social/server/internal/http/middleware/auth"
	corsMiddleware "github.com/vitamin-nn/otus_architect_social/server/internal/http/middleware/cors"
	limitMiddleware "github.com/vitamin-nn/otus_architect_social/server/internal/http/middleware/limit"
	"github.com/vitamin-nn/otus_architect_social/server/internal/http/server"
	"github.com/vitamin-nn/otus_architect_social/server/internal/http/server/messenger/handler"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository/mysql"
)

func messengerCmd(cfg *config.MessengerConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "messenger",
		Short: "starts messenger server",
		Run: func(cmd *cobra.Command, args []string) {
			log.WithFields(cfg.Fields()).Info("starting messenger server")

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
			router := getConfiguredMessengerRouter(messengerRepo, jwt)

			httpSrv := server.NewHTTP(router, cfg.HTTPServer.WriteTimeout, cfg.HTTPServer.ReadTimeout)

			go func() {
				log.Info("starting HTTP messenger server")
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
				log.Error("error while shutdown")
			}

			err := dbShardPool.Close()
			if err != nil {
				log.Fatalf("mysql close connect error: %v", err)
			}
		},
	}
}

func getConfiguredMessengerRouter(messengerRepo repository.MessengerRepo, auth auth.Auth) *gin.Engine {
	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	r.Use(limitMiddleware.MaxAllowed(10))
	r.Use(corsMiddleware.Middleware())

	authorized := r.Group("/api")
	authorized.Use(authMiddleware.Auth(auth))
	{
		authorized.POST("/send", func(c *gin.Context) {
			handler.SendMessage(c, messengerRepo, auth)
		})

		authorized.GET("/dialog/get/:user_id", func(c *gin.Context) {
			handler.GetDialogMessageList(c, messengerRepo, auth)
		})
	}

	return r
}
