package auth

import (
	"CP_Discussion/log"
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := c.Writer
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Status(http.StatusOK)
		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		c.Status(http.StatusOK)
		token := r.Header.Get("Authorization")
		ctx := context.WithValue(r.Context(), "token", token)

		c.Request = r.WithContext(ctx)
		c.Next()
	}
}

func ParseContextToken(ctx context.Context) (string, error) {
	token, ok := ctx.Value("token").(string)
	log.Debug.Println(token)
	if !ok || token == "" || !strings.HasPrefix(token, "Bearer ") {
		return "", errors.New("no token provided")
	}
	bearer := "Bearer "
	token = token[len(bearer):]
	return token, nil
}

func ParseContextClaims(ctx context.Context) (*Claims, error) {
	token, err := ParseContextToken(ctx)
	if err != nil {
		return nil, err
	}
	claims, err := ParseToken(token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
