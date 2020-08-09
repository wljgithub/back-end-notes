package main

import (
	"cookie/handlers"
	"net/http"
)

func main() {
	//cookieName := os.Getenv("COOKIE_NAME")
	//cookieValue := os.Getenv("COOKIE_VALUE")
	//cookieDomain := os.Getenv("COOKIE_DOMAIN")
	//cookiePath := os.Getenv("COOKIE_PATH")
	//cookieDuration := os.Getenv("COOKIE_DURATION")
	//
	//cookieDurationInteger, err := strconv.Atoi(cookieDuration)

	//if err != nil {
	//	panic(err)
	//}
	cookieName := "cookie"
	cookieValue := "value"
	cookieDomain := "localhost"
	//cookiePath := "/"
	cookieDurationInteger := 60

	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("OK"))
	})

	loginHandler := handlers.Login{Name: cookieName,
		Domain: cookieDomain,
		//Path:   cookiePath,
		Value:  cookieValue,
		MaxAge: cookieDurationInteger,
		Next:   handler,
	}

	logoutHandler := handlers.Logout{
		Name:   cookieName,
		Domain: cookieDomain,
		//Path:   cookiePath,
		Next:   handler,
	}

	http.Handle("/login", loginHandler)
	http.Handle("/logout", logoutHandler)
	http.HandleFunc("/hello", handlers.SayHello)

	//address := fmt.Sprintf(":%s", os.Getenv("PORT"))
	address := ":8000"

	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}

