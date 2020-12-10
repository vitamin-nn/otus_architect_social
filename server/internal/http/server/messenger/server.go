package messenger

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vitamin-nn/otus_architect_social/server/internal/auth"
	authMiddleware "github.com/vitamin-nn/otus_architect_social/server/internal/http/middleware/auth"
	corsMiddleware "github.com/vitamin-nn/otus_architect_social/server/internal/http/middleware/cors"
	limitMiddleware "github.com/vitamin-nn/otus_architect_social/server/internal/http/middleware/limit"
	"github.com/vitamin-nn/otus_architect_social/server/internal/http/server/messenger/handler"
	"github.com/vitamin-nn/otus_architect_social/server/internal/repository"
)

type HTTPServerMessenger struct {
	srv *http.Server
}

func NewMessenger(messengerRepo repository.MessengerRepo, auth auth.Auth, wTimeout, rTimeout time.Duration) *HTTPServerMessenger {
	s := new(HTTPServerMessenger)
	router := getConfiguredRouter(messengerRepo, auth)
	s.srv = &http.Server{
		WriteTimeout: wTimeout,
		ReadTimeout:  rTimeout,
		Handler:      router,
	}

	return s
}

func (s *HTTPServerMessenger) Run(addr string) error {
	s.srv.Addr = addr

	return s.srv.ListenAndServe()
}

func (s *HTTPServerMessenger) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func getConfiguredRouter(messengerRepo repository.MessengerRepo, auth auth.Auth) *gin.Engine {
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
