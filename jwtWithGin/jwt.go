package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

var secret = []byte("mysecret")

type Claim struct {
	UserName string
	jwt.StandardClaims
}

func generateToken(para Credential) (string, error) {
	expiry := time.Now().Add(1 * time.Minute).Unix()
	claim := Claim{"jack", jwt.StandardClaims{ExpiresAt: expiry}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(secret)
}

func jwtMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractToken(c)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    1,
				"message": "unauthorized",
			})
			c.Abort()
			return
		}

		err = parseToken(token)
		if err != nil {
			log.Println("failed to parse token")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    1,
				"message": "invalid token",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func extractToken(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")
	var token string
	_, err := fmt.Sscanf(header, "Bearer %s", &token)

	return token, err
}

func parseToken(token string) error {
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	return err
}

func refresh(c *gin.Context) {
	token,err:= extractToken(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "failed to generate token",
		})
		return
	}

	claim:=Claim{}
	jwt.ParseWithClaims(token,&claim, func(token *jwt.Token) (interface{}, error) {
		return secret,nil
	})

	claim.ExpiresAt = time.Now().Add(time.Minute).Unix()

	token ,err = jwt.NewWithClaims(jwt.SigningMethodHS256,claim).SignedString(secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "failed to generate token",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data": struct {
			Token string `json:"token"`
		}{token},
	})
}
