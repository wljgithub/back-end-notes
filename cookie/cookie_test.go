package main

import (
	"bytes"
	"cookie/handlers"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestLogin(t *testing.T) {
	expBody := []byte("test!")

	loginHandler := handlers.Login{
		Name:   "session",
		Value:  "logged in",
		Path:   "/",
		Domain: "test.com",
		MaxAge: 60,
		Next: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Write(expBody)
		}),
	}

	testReq, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("Failed to create test request: %v", err)
	}

	recorder := httptest.NewRecorder()
	loginHandler.ServeHTTP(recorder, testReq)
	response := recorder.Result()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Failed to read the response: %v", err)
	}

	if !bytes.Equal(bodyBytes, expBody) {
		t.Errorf("Unexpected response: %v", err)
	}

	cookies := response.Cookies()

	if len(cookies) != 1 {
		t.Error("Response returned more than one cookie")
	}

	cookie := cookies[0]

	if cookie.Value != "logged in" {
		t.Errorf("Cookie has the wrong value. Expected %v, got %v", "logged in", cookie.Value)
	}

	if cookie.Domain != "test.com" {
		t.Errorf("Cookie has the wrong domain. Expected %v, got %v", "test.com", cookie.Domain)
	}

	if cookie.Name != "session" {
		t.Errorf("Cookie has the wrong name. Expected %v, got %v", "session", cookie.Name)
	}

	if cookie.Path != "/" {
		t.Errorf("Cookie has the wrong domain. Expected %v, got %v", "/", cookie.Path)
	}

	if cookie.MaxAge != 60 {
		t.Errorf("Cookie has the wrong max age. Expected %v, got %v", 60, cookie.MaxAge)
	}
}
func TestLogout(t *testing.T) {
	expBody := []byte("logout")
	logoutHandler := handlers.Logout{
		Name:   "cookie",
		Path:   "/",
		Domain: "localhost",
		Next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(expBody)
		}),
	}

	testReq, err := http.NewRequest(http.MethodGet, "/logout", nil)
	if err != nil {
		t.Errorf("failed to create request: %v", err.Error())
	}

	recorder := httptest.NewRecorder()
	logoutHandler.ServeHTTP(recorder, testReq)
	response := recorder.Result()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("failde to read response body: %v", err.Error())
	}

	assert(t, "Unexpected response:%v", expBody, bodyBytes)

	cookies := response.Cookies()

	assert(t, "Response return more than one cookie", len(cookies), 1)

	cookie := cookies[0]

	assert(t, "Cookie has wrong name,expected: %v got: %v", cookie.Name, logoutHandler.Name)

	if cookie.Value != "" {
		t.Errorf("Cookie has wrong value,expected: %v got: %v", "", cookie.Value)
	}
	assert(t, "Cookie has wrong value,expected", cookie.Value, "")

	assert(t, "Cookie has wrong age,expected", cookie.MaxAge, 0)

}
func assert(t *testing.T, prompt string, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%v got: %v want: %v", prompt, got, want)
	}
}
