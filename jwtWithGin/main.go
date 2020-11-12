package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := gin.Default()
	server.POST("/signIn",signIn)

	auth:=server.Group("",jwtMiddleWare())
	auth.POST("/refresh",refresh)
	auth.GET("ping",ping)
	server.Run()
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK,gin.H{
		"code":0,
		"message":"ok",
	})
}

