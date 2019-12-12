package controllers

import (
	"bytes"
	middlewares "auth-service-template/middlewares"
	"auth-service-template/models"
	"encoding/json"
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

	user := usersTesting[0]
	user.Role = []string{"SUPER"}
	user.CreateOrUpdateWithPlainPw()
	user.Password = usersTesting[0].Password

	_, token, _ := ObtainTokenForTesting(user)
	request, _ := http.NewRequest("GET", "/api/user/list?limit=10&skip=0", nil)
	request.Header.Set("Authorization", token)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")

	userNoAdmin := usersTesting[1]
	userNoAdmin.Role = []string{"NONE"}
	userNoAdmin.CreateOrUpdateWithPlainPw()
	userNoAdmin.Password = usersTesting[1].Password

	_, token, refreshToken := ObtainTokenForTesting(userNoAdmin)
	request, _ = http.NewRequest("GET", "/api/user/list?limit=10&skip=0", nil)
	request.Header.Set("Authorization", token)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, 403, response.Code, "Error expected")

	request.Header.Set("Authorization", refreshToken)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, 401, response.Code, "Error expected")

}

func TestListUsersInOrg(t *testing.T) {
	router := mux.NewRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	router.HandleFunc("/api/user/list_in_organization", ListUsersInOrganization)
	router.Use(middlewares.JwtAuthentication) //attach JWT auth middleware

	user := usersTesting[0]
	user.Role = []string{"ADMIN"}
	user.CreateOrUpdateWithPlainPw()
	user.Password = usersTesting[0].Password
	user.OrganizationPermission[0].OrganziationRole="ADMIN"


	_, token, _:= ObtainTokenForTesting(user)
	request, _ := http.NewRequest("GET", "/api/user/list_in_organization?limit=10&skip=0&organization_id="+ user.OrganizationPermission[0].OrganizationId, nil)
	request.Header.Set("Authorization", token)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")

	request, _ = http.NewRequest("GET", "/api/user/list_in_organization?limit=10&skip=0&organization_id=wrong_organization", nil)
	request.Header.Set("Authorization", token)
	response = httptest.NewRecorder()

	router.ServeHTTP(response, request)
	assert.Equal(t, 403, response.Code, "Error expected")


	userNoAdmin := usersTesting[1]
	userNoAdmin.Role = []string{"NONE"}
	userNoAdmin.CreateOrUpdateWithPlainPw()
	userNoAdmin.Password = usersTesting[1].Password
	userNoAdmin.OrganizationPermission[0].OrganziationRole="NONE"

	_, token, _ = ObtainTokenForTesting(userNoAdmin)
	request, _ = http.NewRequest("GET", "/api/user/list_in_organization?limit=10&skip=0&organization_id="+ userNoAdmin.OrganizationPermission[0].OrganizationId, nil)
	request.Header.Set("Authorization", token)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, 403, response.Code, "Error expected")

}


func TestUpdateOrganizationForUser(t *testing.T) {
	user := usersTesting[0]
	user.OrganizationPermission[0].OrganziationRole= "TEST"
	_, token, _ := ObtainTokenForTesting(user)

	router := mux.NewRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	router.HandleFunc("/api/user/update_organization_permissions", UpdateOrganizationInUser)
	router.Use(middlewares.JwtAuthentication) //attach JWT auth middleware

	bodyPut := user.OrganizationPermission
	jsonPut, _ := json.Marshal(bodyPut)

	request, _ := http.NewRequest("PUT", "/api/user/update_organization_permissions?user_id=" + user.ID, bytes.NewBuffer(jsonPut))
	request.Header.Set("Authorization", token)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	bodyResponse, err := ioutil.ReadAll(response.Body)
	assert.Equal(t, err, nil, "no errors extracting response")

	var responseUser *models.UserResponse
	err = json.Unmarshal(bodyResponse, &responseUser)

	assert.Equal(t, responseUser.Data[0].OrganizationPermission[0] , user.OrganizationPermission[0], "Ok response expected")
	assert.Equal(t, 200, response.Code, "Ok response expected")
}

func TestUpdateOrganizationForUserWithError(t *testing.T) {
	user := usersTesting[3]
	_, token, _:= ObtainTokenForTesting(user)

	router := mux.NewRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	router.HandleFunc("/api/user/update_organization_permissions", UpdateOrganizationInUser)
	router.Use(middlewares.JwtAuthentication) //attach JWT auth middleware

	user.OrganizationPermission[0].OrganizationId= "wrong-organization-test"
	user.OrganizationPermission[0].OrganziationRole= "ADMIN"
	bodyPut := user.OrganizationPermission

	jsonPut, _ := json.Marshal(bodyPut)

	request, _ := http.NewRequest("PUT", "/api/user/update_organization_permissions?user_id=" + user.ID, bytes.NewBuffer(jsonPut))
	request.Header.Set("Authorization", token)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, 403, response.Code, "Ok response expected")
}

func TestUpdateRoleForUser(t *testing.T) {
	user := usersTesting[0]
	user.Role= []string{"SUPER"}
	user.CreateOrUpdateWithPlainPw()
	user.Password = usersTesting[0].Password

	_, token, _:= ObtainTokenForTesting(user)

	router := mux.NewRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	router.HandleFunc("/api/user/update_role", UpdateRoleInUser)
	router.Use(middlewares.JwtAuthentication) //attach JWT auth middleware
	changedRole := "TEST"
	request, _ := http.NewRequest("PUT", "/api/user/update_role?user_id=" + usersTesting[3].ID + "&role="+changedRole, nil)
	request.Header.Set("Authorization", token)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)
	body, err := ioutil.ReadAll(response.Body)

	assert.Equal(t, err, nil, "no errors extracting response")

	var responseUser *models.UserResponse
	err = json.Unmarshal(body, &responseUser)

	assert.Equal(t, 200, response.Code, "Ok response expected")
	assert.Equal(t, responseUser.Data[0].Role[0] , changedRole, "Ok response expected")

}


func TestUpdateRoleForUserWithError(t *testing.T) {
	user := usersTesting[3]
	user.Role= []string{"NONE"}
	user.CreateOrUpdateWithPlainPw()
	user.Password = usersTesting[3].Password

	_, token, _:= ObtainTokenForTesting(user)

	router := mux.NewRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	router.HandleFunc("/api/user/update_role", UpdateRoleInUser)
	router.Use(middlewares.JwtAuthentication) //attach JWT auth middleware
	changedRole := "TEST"
	request, _ := http.NewRequest("PUT", "/api/user/update_role?user_id=" + usersTesting[3].ID + "&role="+changedRole, nil)
	request.Header.Set("Authorization", token)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)
	assert.Equal(t, 403, response.Code, "Error expected")

}

func ObtainTokenForTesting(user models.User) (err error, token string, refreshToken string) {
	_ = SignInHelperMux()
	loginResponse := LogInHelperWithUser(user)
	body, _ := ioutil.ReadAll(loginResponse.Body)
	bodyMap := models.UserResponse{}
	err = json.Unmarshal(body, &bodyMap)

	token = "Bearer " + bodyMap.Data[0].Token
	refreshToken = "Bearer " + bodyMap.Data[0].RefreshToken

	return err, token, refreshToken
}
