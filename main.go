package main

import (
	"arch/internal/config"
	"arch/internal/db"
	"arch/internal/handler"
	"arch/internal/server"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

func main() {

	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	gormDB, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(gormDB)

	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()

	router.Get("/ping", handler.Ping)

	srv := server.NewServer("8080", router)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Произошла ошибка", err)
		}
	}()

	<-ctx.Done()

	shutDownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	if err := srv.HttpServer.Shutdown(shutDownCtx); err != nil {
		fmt.Println("Ошибка при shutDown", err)
	}

	fmt.Println("Остановили сервер")

}
