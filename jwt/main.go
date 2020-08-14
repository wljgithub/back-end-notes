package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"time"
)

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}
var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Password string `json:"password"`
	UserName string `json:"user_name"`
}

type Claims struct {
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

func Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expectedPassword, ok := users[creds.UserName]

	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		UserName: creds.UserName,
		StandardClaims: jwt.StandardClaims{
			// in jwt,expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokenStr := c.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Fprint(w, fmt.Sprintf("Welcome %s", claims.UserName))
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokenStr := c.Value
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Minute {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	twt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := twt.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	fmt.Fprint(w,"refresh jwt successfully")

}
func main() {

	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/welcome", Welcome)
	http.HandleFunc("/refresh", Refresh)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
