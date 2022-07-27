package fileHandler

import (
	"CP_Discussion/file/fileManager"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		year, err := strconv.ParseInt(c.Param("year"), 10, 32)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
		}
		semester, err := strconv.ParseInt(c.Param("semester"), 10, 32)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
		}
		filename := c.Param("filename")
		mdPath := fileManager.BuildPostPath(int(year), int(semester), filename)
		c.File(mdPath)
	}
}
func AvatarHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Param("filename")
		avatarPath := fileManager.BuildAvatarPath(filename)
		c.File(avatarPath)
	}
}
