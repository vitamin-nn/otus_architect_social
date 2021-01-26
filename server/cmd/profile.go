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
	"github.com/vitamin-nn/otus_architect_social/server/internal/cache"
	cacheRedis "github.com/vitamin-nn/otus_architect_social/server/internal/cache/redis"
	"github.com/vitamin-nn/otus_architect_social/server/internal/config"
	authMiddleware "github.com/vitamin-nn/otus_architect_social/server/internal/http/middleware/auth"
	corsMiddleware "github.com/vitamin-nn/otus_architect_social/server/internal/http/middleware/cors"
	limitMiddleware "github.com/vitamin-nn/otus_architect_social/server/internal/http/middleware/limit"
	"github.com/vitamin-nn/otus_architect_social/server/internal/http/server"
	"github.com/vitamin-nn/otus_architect_social/server/internal/http/server/profile/handler"
	"github.com/vitamin-nn/otus_architect_social/server/internal/queue/rabbit"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository/mysql"
)

func serverCmd(cfg *config.SocialConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "profile",
		Short: "Starts social network server",
		Run: func(cmd *cobra.Command, args []string) {
			log.WithFields(cfg.Fields()).Info("starting social server")

			jwt := jwtAuth.New(cfg.JWT.Secret, cfg.JWT.AccessLifeTime, cfg.JWT.RefreshLifeTime)

			dbPool := NewDBPool(context.Background(), cfg.MySQL.GetDSN(), cfg.MySQL.SlavesDSN)

			profileRepo := mysql.NewProfileRepo(dbPool)
			feedRepo := mysql.NewFeedRepo(dbPool)

			redisPool := cacheRedis.NewFeedCache(newRedisPool(context.Background(), cfg.Redis.Addr), maxFeedLen)

			rmq := rabbit.NewPublish(cfg.Rabbit.GetAddr(), cfg.Rabbit.ExchangeName, cfg.Rabbit.QueueName)
			err := rmq.Connect()
			if err != nil {
				log.Fatalf("rabbit connect error: %v", err)
			}
			defer rmq.Close()

			router := getConfiguredProfileRouter(profileRepo, feedRepo, jwt, redisPool, rmq)
			httpSrv := server.NewHTTP(router, cfg.HTTPServer.WriteTimeout, cfg.HTTPServer.ReadTimeout)

			go func() {
				log.Info("starting HTTP server")
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
				log.Error("error while shutdown")
			}

			err = dbPool.Close()
			if err != nil {
				log.Fatalf("mysql close connect error: %v", err)
			}
		},
	}
}

//nolint:funlen
func getConfiguredProfileRouter(profileRepo repository.ProfileRepo, feedRepo repository.FeedRepo, auth auth.Auth, cacheFeed cache.Feed, rmq *rabbit.Publish) *gin.Engine {
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
		authorized.GET("/", func(c *gin.Context) {
			handler.Main(c, profileRepo)
		})

		authorized.POST("/login", func(c *gin.Context) {
			handler.Login(c, profileRepo, auth)
		})

		authorized.GET("/user/:id", func(c *gin.Context) {
			handler.PublicProfile(c, profileRepo)
		})

		authorized.GET("/profile", func(c *gin.Context) {
			handler.MyProfile(c, profileRepo, auth)
		})

		authorized.GET("/profile/filter", func(c *gin.Context) {
			handler.FilteredProfile(c, profileRepo)
		})

		authorized.GET("/friends", func(c *gin.Context) {
			handler.FriendList(c, profileRepo, auth)
		})

		authorized.POST("/friends/add", func(c *gin.Context) {
			handler.AddFriend(c, profileRepo, auth)
		})

		authorized.POST("/friends/remove", func(c *gin.Context) {
			handler.RemoveFriend(c, profileRepo, auth)
		})

		authorized.POST("/register", func(c *gin.Context) {
			handler.Register(c, profileRepo, auth)
		})

		authorized.POST("/feed/add", func(c *gin.Context) {
			handler.AddFeedMessage(c, feedRepo, auth, rmq)
		})

		authorized.GET("/feed/get/:id", func(c *gin.Context) {
			handler.GetFeed(c, cacheFeed)
		})
	}

	return r
}
