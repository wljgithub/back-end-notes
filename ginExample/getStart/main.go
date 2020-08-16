package main

import (
	"back-end-notes/ginExample/getStart/ginExample/getStart/controller"
	"back-end-notes/ginExample/getStart/ginExample/getStart/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	videoService    = service.New()
	videoController = controller.New(videoService)
)

func main() {
	server := gin.Default()

	server.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK!",
		})
	})
	server.GET("/videos", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, videoController.FindAll())
	})
	server.POST("/videos", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, videoController.Save(ctx))
	})
	server.Run(":8080")
}
