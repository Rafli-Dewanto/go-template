package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/Rafli-Dewanto/go-template/internal/context"
	"github.com/Rafli-Dewanto/go-template/internal/utils"
)

type Middleware func(http.Handler) http.Handler

// Chain applies middlewares to a http.Handler in the order they are passed
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

// Logger logs the incoming HTTP request and its duration
func Logger(logger *utils.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			logger.Info("Started %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)

			logger.Info("Completed in %v", time.Since(start))
		})
	}
}

// Recover recovers from panics and logs the error
func Recover(logger *utils.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					stack := debug.Stack()
					logger.Error("PANIC: %v\n%s", err, string(stack))
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// CORS handles Cross-Origin Resource Sharing
func CORS() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequestID adds a unique request ID to each request
func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := fmt.Sprintf("%d", time.Now().UnixNano())
			w.Header().Set("X-Request-ID", requestID)
			next.ServeHTTP(w, r)
		})
	}
}

// APIID adds a unique API ID to each request for tracking
func APIID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if API ID already exists in request header
			apiID := r.Header.Get("X-API-ID")
			if apiID == "" {
				// Generate new API ID using crypto utility
				apiID = utils.GenerateAPIID()
			}

			// Add API ID to response header
			w.Header().Set("X-API-ID", apiID)

			// Add API ID to request context
			ctx := context.WithAPIID(r.Context(), apiID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
