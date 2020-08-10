package handlers

import (
	scookie "github.com/gorilla/securecookie"
	"log"
	"net/http"
)

// Login is a handler that sets a logged in cookie
type Login struct {
	Name   string
	Value  string
	Path   string
	Domain string
	MaxAge int
	Next   http.Handler
}

var hashkey = scookie.GenerateRandomKey(64)
var blockkey = scookie.GenerateRandomKey(32)
var s = scookie.New(hashkey,blockkey)

func (l Login) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if encode ,err := s.Encode(l.Name,l.Value);err==nil{
		cookie := http.Cookie{
			Name:     l.Name,
			Value:    encode,
			Domain:   l.Domain,
			Path:     l.Path,
			MaxAge:   l.MaxAge,
			HttpOnly: true,
		}
		http.SetCookie(rw, &cookie)
		l.Next.ServeHTTP(rw, r)
		return
	}
	log.Println("failed to set cookie")
	// todo: otherwise
}