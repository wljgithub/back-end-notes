package handlers

import "net/http"

func SayHello(w http.ResponseWriter,r *http.Request){
	w.Write([]byte("hello"))
}
