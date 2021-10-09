package main

import (
	"Instagram-API/routing"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Router() *routing.Router {
	mux := routing.NewRouter()
	mux.HandleFunc("/users", getUsers).Methods("GET")
	mux.HandleFunc("/users/616159235ccd11baa4ec2d1c", getUser).Methods("GET")
	mux.HandleFunc("/posts/6161b3e37430caa4e1ba575d", getPost).Methods("GET")
	mux.HandleFunc("/posts/users/616159235ccd11baa4ec2d1c", getPosts).Methods("GET")
	mux.HandleFunc("/posts", getPosts).Methods("POST")
	return mux
}
func TestGetUsers(t *testing.T) {
	request, _ := http.NewRequest("GET", "/users", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	res := response.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Experted status ok; got %v", res.StatusCode)
	}
	defer res.Body.Close()

	_, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}
}
func TestGetUser(t *testing.T) {
	request, _ := http.NewRequest("GET", "/users", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	res := response.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Experted status ok; got %v", res.StatusCode)
	}
	defer res.Body.Close()

	_, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}
}
func TestGetPost(t *testing.T) {
	request, _ := http.NewRequest("GET", "/posts/6161b3e37430caa4e1ba575d", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	res := response.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Experted status ok; got %v", res.StatusCode)
	}
	defer res.Body.Close()

	_, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}
}
func TestGetPosts(t *testing.T) {
	request, _ := http.NewRequest("GET", "/posts/users/616159235ccd11baa4ec2d1c", nil)
	response := httptest.NewRecorder()

	Router().ServeHTTP(response, request)
	res := response.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Experted status ok; got %v", res.StatusCode)
	}
	defer res.Body.Close()

	_, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}
}
func TestCreatePost(t *testing.T) {
	request, _ := http.NewRequest("GET", "/posts", nil)
	response := httptest.NewRecorder()

	Router().ServeHTTP(response, request)
	res := response.Result()
	if res.StatusCode != 405 {
		t.Errorf("Experted status ok; got %v", res.StatusCode)
	}
	defer res.Body.Close()

	_, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}
}
func TestCreateUser(t *testing.T) {
	request, _ := http.NewRequest("GET", "/posts", nil)
	response := httptest.NewRecorder()

	Router().ServeHTTP(response, request)
	res := response.Result()
	if res.StatusCode != 405 {
		t.Errorf("Experted status ok; got %v", res.StatusCode)
	}
	defer res.Body.Close()

	_, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}
}
