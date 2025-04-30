package middleware

import (
	"net/http"
	"strings"

	"github.com/Rafli-Dewanto/go-template/internal/auth"
	"github.com/Rafli-Dewanto/go-template/internal/handler"
)

func AuthMiddleware(tokenManager *auth.TokenManager) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				handler.WriteErrorResponse(w, http.StatusUnauthorized, "Missing authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				handler.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			tokenString := parts[1]
			claims, err := tokenManager.ValidateToken(tokenString)
			if err != nil {
				switch err {
				case auth.ErrExpiredToken:
					handler.WriteErrorResponse(w, http.StatusUnauthorized, "Token has expired")
				default:
					handler.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid token")
				}
				return
			}

			// Add user claims to request context
			ctx := r.Context()
			ctx = auth.WithUserClaims(ctx, claims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
