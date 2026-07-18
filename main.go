package main

import (
	// "arch/internal/cache"
	"arch/internal/db"
	"arch/internal/handler"
	"arch/internal/repository"
	"arch/internal/server"
	"arch/internal/services"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/patrickmn/go-cache"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

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

	userRepo := &repository.UserRepo{DB: gormDB, Cache: cacheUser}
	userService := &services.UserService{UserRep: userRepo}
	userHandler := &handler.UserHandler{UserServices: userService}

	r := chi.NewRouter()
	userHandler.RegisterRoutes(r)

	srv := server.NewServer("8080", r)

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
