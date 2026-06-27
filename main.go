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

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	gormDB, err := db.FirstPostgresLoad(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(gormDB)

	router := chi.NewRouter()

	h := handler.NewHandler(gormDB)

	router.Get("/ping", handler.Ping)
	router.Get("/get-db", h.HandlerApp)

	srv := server.NewServer("8080", router)

	// создали контекст, который ловит сигналы и корректное завершает программу
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// запускаем сервер
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Произошла ошибка", err)
		}
	}()

	// ждем сигнали о завершении
	<-ctx.Done()

	//Даем серверу 3 секунд на завершение текущих запросов.
	shutDownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := srv.HttpServer.Shutdown(shutDownCtx); err != nil {
		fmt.Println("Ошибка при shutDown", err)
	}

	fmt.Println("Остановили сервер")

}
