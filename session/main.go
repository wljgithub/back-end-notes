package main

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"net/http"
)

var cache redis.Conn

func main() {
	initCache()

	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/welcome", Welcome)
	http.HandleFunc("/refresh",Refresh)
	log.Println("server is running on: 8000")
	log.Fatalln(http.ListenAndServe(":8000", nil))
}


func initCache() {
	conn, err := redis.DialURL("redis://localhost")
	if err != nil {
		log.Fatalln("failed to init redis:", err)
	}
	cache = conn

}
