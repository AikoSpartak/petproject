package app

import (
	"arch/internal/handler"
	"arch/internal/handler/HandlerOrders"
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

	OrderRepo    *repository.OrderRepositoryStruct
	OrderService *services.OrdersServicesStruct
	OrderHandler *HandlerOrders.OrderHandler

	Router *chi.Mux
}

func AppService(db *gorm.DB, cache *cache.Cache) *App {

	userRepo := &repository.UserRepositor{DB: db, Cache: cache}
	userService := &services.UserService{UserRepositor: userRepo}
	userHandler := &handler.UserHandler{UserServices: userService}

	orderRepo := &repository.OrderRepositoryStruct{DB: db, Cache: cache}
	orderService := &services.OrdersServicesStruct{OrdersResotitory: orderRepo}
	orderHandler := &HandlerOrders.OrderHandler{OrderServices: orderService}

	router := chi.NewRouter()

	userHandler.RegisterRoutes(router)
	orderHandler.RegisterRouts(router)

	return &App{
		UserRepo:     userRepo,
		UserService:  userService,
		UserHandler:  userHandler,
		OrderHandler: orderHandler,
		Router:       router,
	}
}
