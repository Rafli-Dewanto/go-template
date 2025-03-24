package router

import (
	"net/http"

	"github.com/Rafli-Dewanto/go-template/internal/handler"
	"github.com/Rafli-Dewanto/go-template/internal/repository"
	"github.com/Rafli-Dewanto/go-template/internal/service"
	"github.com/Rafli-Dewanto/go-template/internal/utils"
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
	mux := http.NewServeMux()

	// User routes
	mux.HandleFunc("/users", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			r.userHandler.Create(w, req)
		case http.MethodGet:
			r.userHandler.List(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/users/", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			r.userHandler.GetByID(w, req)
		case http.MethodPut:
			r.userHandler.Update(w, req)
		case http.MethodPatch:
			r.userHandler.SoftDelete(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
