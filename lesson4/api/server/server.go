package server

import (
	"context"
	"gb-backend/lesson4/app/repos/files"
	"gb-backend/lesson4/app/starter"
	"net/http"
	"time"
)

var _ starter.APIServer = &Server{}

type Server struct {
	srv http.Server
	f   *files.Files
}

func NewServer(addr string, h http.Handler) *Server {
	s := &Server{}

	s.srv = http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
	}
	return s
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	s.srv.Shutdown(ctx)
	cancel()
}

func (s *Server) Start(f *files.Files) {
	s.f = f
	go s.srv.ListenAndServe()
}
