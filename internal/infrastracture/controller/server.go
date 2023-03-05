package controller

import (
	"context"
	"net/http"
	"time"
)

type HTTPServer struct {
	httpServer *http.Server
}

func NewHTTPServer(port string, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		httpServer: &http.Server{
			Addr:           ":" + port,
			Handler:        handler,
			MaxHeaderBytes: 1 << 20,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
		},
	}
}

func (s *HTTPServer) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
