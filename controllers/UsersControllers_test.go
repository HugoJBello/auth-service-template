package controllers

import (
	middlewares "auth-service-template/middlewares"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"io/ioutil"

	mux "github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestListUsers(t *testing.T) {
	router := mux.NewRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	router.HandleFunc("/api/user/list", ListUsers)
	router.Use(middlewares.JwtAuthentication) //attach JWT auth middleware

	_, token := ObtainTokenForTesting()

	request, _ := http.NewRequest("GET", "/api/user/list?limit=1&skip=0", nil)
	request.Header.Set("Authorization", token)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")
	fmt.Println(response.Body)

}

func ObtainTokenForTesting() (error, string) {
	loginResponse := LogInHelperMux()
	body, _ := ioutil.ReadAll(loginResponse.Body)
	bodyMap := make(map[string]map[string]string)
	err := json.Unmarshal(body, &bodyMap)

	token := "Bearer " + bodyMap["user"]["token"]

	return err, token
}
