package fileHandler

import (
	"CP_Discussion/log"
	"path"

	"github.com/gin-gonic/gin"
)

const (
	PostPath   = "data/post"
	AvatarPath = "data/avatar"
)

func PostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		year := c.Param("year")
		semester := c.Param("semester")
		filename := c.Param("id") + ".md"
		dir := path.Join(PostPath, year, semester)
		c.File(path.Join(dir, filename))
	}
}
func AvatarHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Param("filename")
		dir := path.Join(AvatarPath)
		log.Debug.Printf("%s\n", dir)
		c.File(path.Join(dir, filename))
	}
}
