package server

import (
	"net/http"
)

type Server struct {
	HttpServer *http.Server //
}

func NewServer(port string, handler http.Handler) *Server {
	return &Server{
		HttpServer: &http.Server{
			Addr:    ":" + port,
			Handler: handler,
		},
	}
}

func (s *Server) ListenAndServe() error {
	return s.HttpServer.ListenAndServe()
}
