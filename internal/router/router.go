package router

import (
	"net/http"

	"github.com/Rafli-Dewanto/go-template/internal/handler"
	"github.com/Rafli-Dewanto/go-template/internal/repository"
	"github.com/Rafli-Dewanto/go-template/internal/service"
	"github.com/Rafli-Dewanto/go-template/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	customMiddleware "github.com/Rafli-Dewanto/go-template/internal/middleware"
	"github.com/jmoiron/sqlx"
)

type Router struct {
	userHandler *handler.UserHandler
}

func NewRouter(db *sqlx.DB) *Router {
	// logger
	logger, err := utils.NewLogger("files/log/app.log")
	if err != nil {
		panic(err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db, logger)

	// Initialize services
	userService := service.NewUserService(userRepo, logger)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService, logger)

	return &Router{
		userHandler: userHandler,
	}
}

func (r *Router) SetupRoutes() http.Handler {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(customMiddleware.APIID())
	router.Use(customMiddleware.CORS())

	// User routes
	router.Route("/users", func(route chi.Router) {
		route.Get("/", r.userHandler.List)
		route.Post("/", r.userHandler.Create)
		route.Get("/{id}", r.userHandler.GetByID)
		route.Put("/{id}", r.userHandler.Update)
		route.Patch("/{id}", r.userHandler.SoftDelete)
	})

	return router
}
