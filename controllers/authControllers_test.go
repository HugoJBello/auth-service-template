package controllers

import (
	"auth-service-template/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	response := SignInHelperMux()
	assert.Equal(t, 200, response.Code, "OK response is expected")

}

func TestLogIn(t *testing.T) {
	response := LogInHelperMux()
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func SignInHelperMux() (response *httptest.ResponseRecorder) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/api/user/new", CreateUser)
	body := models.User{Username: "testAdmin", Password: "password2", Role: []string{"ADMIN"}}
	json, _ := json.Marshal(body)

	request, _ := http.NewRequest("POST", "/api/user/new", bytes.NewBuffer(json))
	response = httptest.NewRecorder()

	mux.ServeHTTP(response, request)
	return response
}

func LogInHelperMux() (response *httptest.ResponseRecorder) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/api/user/login", Authenticate)
	body := models.User{Username: "testAdmin", Password: "password2"}
	json, _ := json.Marshal(body)

	request, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(json))
	response = httptest.NewRecorder()

	mux.ServeHTTP(response, request)
	return response
}
