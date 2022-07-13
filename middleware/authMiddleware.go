package middleware

import (
	"CP_Discussion/auth"
	"context"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}
		bearer := "Bearer "
		token = token[len(bearer):]

		claims, err := auth.ParseToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), string("auth"), claims)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
