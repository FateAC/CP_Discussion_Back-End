package fileHandler

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func FileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.File(filepath.Join("data", c.Request.URL.String()))
	}
}
