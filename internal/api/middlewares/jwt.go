package middlewares

import (
	"context"
	"net/http"
	"restproject/internal/auth"
)

type ContextKey string

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("Bearer")
		if err != nil {
			http.Error(w, "authorization header missing", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ParseToken(token.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), auth.ContextKeyRole, claims["role"])
		ctx = context.WithValue(ctx, auth.ContextKeyExpiresAt, claims["exp"])
		ctx = context.WithValue(ctx, auth.ContextKeyUsername, claims["user"])
		ctx = context.WithValue(ctx, auth.ContextKeyUserID, claims["uid"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
