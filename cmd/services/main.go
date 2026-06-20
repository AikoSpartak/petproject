package main

import (
	"arch/cmd/internal/handler"
	"arch/cmd/internal/server"

	"github.com/go-chi/chi"
)

func main() {

	router := chi.NewRouter()

	router.Get("/ping", handler.Ping)

	srv := &server.Server{}
	if err := srv.ServerRun("8080", router); err != nil {
		panic(err)
	}

}
