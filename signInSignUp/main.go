package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type Credential struct {
	UserName string `json:"user_name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

var credentials = make(map[string]Credential)

func signUp(c *gin.Context) {
	var req Credential
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "bind parameter error",
		})
		return
	}

	if _, ok := credentials[req.UserName]; ok {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "user has register",
		})
		return
	}

	// encrypt user password
	hashByte, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("failed to hash password")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "failed to hash password"})
		return
	}

	req.Password = string(hashByte)
	credentials[req.UserName] = req

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

func signIn(c *gin.Context) {
	var req Credential
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "bind parameter error",
		})
		return
	}

	credential, ok := credentials[req.UserName]
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "user hasn't sign up",
		})
		return
	}

	// compare the password that user upload and the password that store in server
	if err := bcrypt.CompareHashAndPassword([]byte(credential.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "incorrect password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
	})

}

func main() {

	server := gin.Default()
	server.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{"message":"hello"})
	})
	server.POST("/signIn", signIn)
	server.POST("/signUp", signUp)

	server.Run(":8080")
}
