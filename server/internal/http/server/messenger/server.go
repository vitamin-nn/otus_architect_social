package messenger

import (
	"context"
	"net/http"
	"time"
)

type HTTPServerMessenger struct {
	srv *http.Server
}

func NewMessenger(router http.Handler, wTimeout, rTimeout time.Duration) *HTTPServerMessenger {
	s := new(HTTPServerMessenger)

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