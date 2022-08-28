package auth

import (
	"CP_Discussion/log"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

func RefreshHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		claims, err := ParseContextClaims(ctx)
		if err != nil {
			log.Debug.Print(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		userID := claims.UserID
		refreshToken := ctx.Value("token").(string)
		accessToken, err := CreateToken(time.Now(), time.Now(), time.Now().Add(time.Hour), userID)
		if err != nil {
			err = errors.Wrap(err, "create access token failed")
			log.Debug.Print(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Debug.Printf("create access token for userID %s: %s", userID, accessToken)
		if time.Now().Add(time.Duration(4) * time.Hour).After(claims.ExpiresAt.Time) {
			refreshToken, err = CreateToken(time.Now(), time.Now(), time.Now().Add(time.Duration(24)*time.Hour), userID)
			if err != nil {
				err = errors.Wrap(err, "create refresh token failed")
				log.Debug.Print(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			log.Debug.Printf("create refresh token for userID %s: %s", userID, refreshToken)
		}
		c.JSON(http.StatusOK, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}
