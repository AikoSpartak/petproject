package main

import (
	// "arch/internal/cache"
	"arch/internal/app"
	"arch/internal/db"
	"arch/internal/logger"
	"arch/internal/server"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/patrickmn/go-cache"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	logger.InitLogger()

	gormDB, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal(err)
	}

	defer sqlDB.Close()

	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	cacheUser := cache.New(5*time.Minute, 10*time.Minute)

	r := app.AppService(gormDB, cacheUser)
	srv := server.NewServer("8080", r.Router)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Произошла ошибка", err)
		}
	}()

	<-ctx.Done()

	slog.Info("получен сигнал о щавершинии")

	stop()

	shutDownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.HttpServer.Shutdown(shutDownCtx); err != nil {
		fmt.Println("Ошибка при shutDown", err)
	}

	fmt.Println("Остановили сервер")

}
