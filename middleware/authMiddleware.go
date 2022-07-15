package middleware

import (
	"CP_Discussion/auth"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		w, r := c.Writer, c.Request
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

		token := r.Header.Get("Authorization")
		if token == "" {
			c.Status(http.StatusOK)
			c.Next()
			return
		}
		bearer := "Bearer "
		token = token[len(bearer):]

		claims, err := auth.ParseToken(token)
		if err != nil {
			c.Status(http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), string("auth"), claims)

		c.Request = r.WithContext(ctx)
		c.Next()
	}
}
