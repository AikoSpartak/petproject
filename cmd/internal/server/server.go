package server

import (
	"fmt"
	"net/http"
)

type Server struct {
	httpServer http.Server
}

func (s *Server) ServerRun(port string, handler http.Handler) error {

	s.httpServer = http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	fmt.Println("Сервер запустился")

	return s.httpServer.ListenAndServe()
}
