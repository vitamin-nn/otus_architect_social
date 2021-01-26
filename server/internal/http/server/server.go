package server

import (
	"context"
	"net/http"
	"time"
)

type HTTP struct {
	srv *http.Server
}

func NewHTTP(router http.Handler, wTimeout, rTimeout time.Duration) *HTTP {
	s := new(HTTP)
	s.srv = &http.Server{
		WriteTimeout: wTimeout,
		ReadTimeout:  rTimeout,
		Handler:      router,
	}

	return s
}

func (s *HTTP) Run(addr string) error {
	s.srv.Addr = addr

	return s.srv.ListenAndServe()
}

func (s *HTTP) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
