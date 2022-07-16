package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
		w, r := c.Writer, c.Request
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

		token := r.Header.Get("Authorization")
		ctx := context.WithValue(r.Context(), string("token"), token)

		c.Request = r.WithContext(ctx)
		c.Next()
	}
}
