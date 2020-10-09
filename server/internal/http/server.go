package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vitamin-nn/otus_architect_social/server/internal/auth"
	"github.com/vitamin-nn/otus_architect_social/server/internal/http/handler"
	authMiddleware "github.com/vitamin-nn/otus_architect_social/server/internal/http/middleware/auth"
	corsMiddleware "github.com/vitamin-nn/otus_architect_social/server/internal/http/middleware/cors"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository"
)

type HTTPServer struct {
	srv *http.Server
}

func New(profileRepo repository.ProfileRepo, auth auth.Auth, wTimeout, rTimeout time.Duration) *HTTPServer {
	s := new(HTTPServer)
	router := getConfiguredRouter(profileRepo, auth)
	s.srv = &http.Server{
		WriteTimeout: wTimeout,
		ReadTimeout:  rTimeout,
		Handler:      router,
	}
	return s
}

func (s *HTTPServer) Run(addr string) error {
	s.srv.Addr = addr
	return s.srv.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func getConfiguredRouter(profileRepo repository.ProfileRepo, auth auth.Auth) *gin.Engine {
	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	r.Use(corsMiddleware.CORSMiddleware())

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
	}

	return r
}
