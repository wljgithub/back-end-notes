package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var credential = map[string]string{
	"username": "jack",
	"password":"password",
}

type Credential struct {
	UserName string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}


func signIn(c *gin.Context) {
	var req Credential
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("bind para error")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "bind para error",
		})
		return
	}

	if req.UserName != credential["username"]|| req.Password != credential["password"]{
		c.JSON(http.StatusUnauthorized,gin.H{
			"code":1,
			"message":"incorrect username or password",
		})
		return
	}

	// generate jwt token
	token ,err:= generateToken(req)
	if err!= nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"code":1,
			"message":"failed to generate token",
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"code":0,
		"message":"ok",
		"data": struct {
			Token string `json:"token"`
		}{token},
	})

}
