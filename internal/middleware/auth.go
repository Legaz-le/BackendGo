package middleware

import (
	"context"
	"net/http"
	"strings"

	"example.com/mod/internal/auth"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		bearerCheck := strings.HasPrefix(header, "Bearer ")

		if !bearerCheck {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")

		check, err := auth.ValidateAccessToken(token)

		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsKey, check)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
