package app

import (
	"arch/internal/handler"
	"arch/internal/repository"
	"arch/internal/services"

	"github.com/patrickmn/go-cache"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

type App struct {
	UserRepo    *repository.UserRepositor
	UserService *services.UserService
	UserHandler *handler.UserHandler
	Router      *chi.Mux
}

func AppService(db *gorm.DB, cache *cache.Cache) *App {

	userRepo := &repository.UserRepositor{DB: db, Cache: cache}
	userService := &services.UserService{UserRepositor: userRepo}
	userHandler := &handler.UserHandler{UserServices: userService}

	router := chi.NewRouter()

	userHandler.RegisterRoutes(router)

	return &App{
		UserRepo:    userRepo,
		UserService: userService,
		UserHandler: userHandler,
		Router:      router,
	}
}
