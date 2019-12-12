package controllers

import (
	"bytes"
	"auth-service-template/middlewares"
	"auth-service-template/models"
	"auth-service-template/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var usersTesting, _ = utils.GetUsersForTesting()

func TestCreateUser(t *testing.T) {
	response := SignInHelperMux()
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestBadLogin(t *testing.T) {
	user := usersTesting[0]
	user.Password = user.Password + "-"
	response := LogInHelperWithUser(user)
	assert.Equal(t, 401, response.Code, "error is expected")
}

func TestLogIn(t *testing.T) {
	user := usersTesting[0]
	response := LogInHelperWithUser(user)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestRefresh(t *testing.T) {
	user := usersTesting[0]
	response := LogInHelperWithUser(user)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func SignInHelperMux() (response *httptest.ResponseRecorder) {
	user := usersTesting[0]
	return SignInHelperWithUser(user)
}

func LogInHelperMux() (response *httptest.ResponseRecorder) {
	user := usersTesting[0]
	return LogInHelperWithUser(user)
}

func LogInHelperWithUser(user models.User) (response *httptest.ResponseRecorder) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/api/user/login", Authenticate)
	body := user
	json, _ := json.Marshal(body)

	request, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(json))
	response = httptest.NewRecorder()

	mux.ServeHTTP(response, request)
	return response
}

func SignInHelperWithUser(user models.User) (response *httptest.ResponseRecorder) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/api/user/new", CreateUser)
	body := user
	json, _ := json.Marshal(body)

	request, _ := http.NewRequest("POST", "/api/user/new", bytes.NewBuffer(json))
	response = httptest.NewRecorder()

	mux.ServeHTTP(response, request)
	return response
}

func TestRefreshToken(t *testing.T) {
	user := usersTesting[1]

	_, _, refreshToken:= ObtainTokenForTesting(user)

	router := mux.NewRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	router.HandleFunc("/api/user/refresh", Refresh)
	router.Use(middlewares.JwtAuthentication) //attach JWT auth middleware

	request, _ := http.NewRequest("POST", "/api/user/refresh", nil)
	request.Header.Set("Authorization", refreshToken)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK")

	body, err := ioutil.ReadAll(response.Body)

	assert.Equal(t, err, nil, "no errors extracting response")
	var responseRefresh *models.RefreshResponse
	err = json.Unmarshal(body, &responseRefresh)

	assert.NotEqual(t, responseRefresh.Token, "", "should retrieve token")
	assert.NotEqual(t, responseRefresh.RefreshToken, "", "should retrieve refresh token")
}
