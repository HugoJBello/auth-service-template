package controllers

import (
	"auth-service-template/models"
	"auth-service-template/utils"
	"encoding/json"
	"github.com/gorilla/context"
	"net/http"
	"strconv"
)

var ListUsers = func(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	skip, _ := strconv.Atoi(r.FormValue("skip"))
	var role  = context.Get(r, "role").([]string)
	if utils.Contains(role, "SUPER"){
		users := models.GetUsers(limit, skip)
		response := models.UserResponse{Code: 200, Message: "users retrieved", Data: users}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		response := models.UserResponse{Code: 404, Message: "not allowed", Data: nil}
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}

}

var ListUsersInOrganization = func(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	skip, _ := strconv.Atoi(r.FormValue("skip"))
	organizationId := r.FormValue("organization_id")

	var role  = context.Get(r, "role").([]string)
	var organizationIdContext  = context.Get(r, "organization_permissions").(map[string]string)


	if utils.Contains(role, "SUPER") || organizationIdContext[organizationId] == "ADMIN" {
		users := models.GetUsersInOrg(limit, skip, organizationId)
		response := models.UserResponse{Code: 200, Message: "users retrieved", Data: users}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		response := models.UserResponse{Code: 404, Message: "not allowed", Data: nil}
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}

}

var allowedForOrganizations = func(permissions []models.OrganizationPermission, permissionsMapInContext map[string]string) bool{
	for _, organization := range(permissions){
		organizationId := organization.OrganizationId
		if (permissionsMapInContext[organizationId]!="ADMIN"){
			return false
		}
	}
	return true
}

var UpdateOrganizationInUser = func(w http.ResponseWriter, r *http.Request) {

	organizationPermisions := &[]models.OrganizationPermission{}
	err := json.NewDecoder(r.Body).Decode(organizationPermisions)

	if err != nil {
		response := models.UserResponse{Code: 401, Message: "wrong format", Data: nil}
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}

	var organizationIdContext  = context.Get(r, "organization_permissions").(map[string]string)

	userId := r.FormValue("user_id")
	var role  = context.Get(r, "role").([]string)
	allowed := allowedForOrganizations(*organizationPermisions, organizationIdContext)
	if utils.Contains(role, "SUPER") || allowed {
		models.UpdateOrganizationInUser(*organizationPermisions ,userId)
		user := models.GetUserById(userId)
		response := models.UserResponse{Code: 200, Message: "users retrieved", Data: []models.User{*user}}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		response := models.UserResponse{Code: 404, Message: "not allowed", Data: nil}
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}


var UpdateRoleInUser = func(w http.ResponseWriter, r *http.Request) {

	roleInput := []string{r.FormValue("role")}
	userId := r.FormValue("user_id")

	var role  = context.Get(r, "role").([]string)
	if utils.Contains(role, "SUPER") {
		models.UpdateRoleInUser(roleInput ,userId)
		user := models.GetUserById(userId)
		response := models.UserResponse{Code: 200, Message: "users retrieved", Data: []models.User{*user}}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		response := models.UserResponse{Code: 404, Message: "not allowed", Data: nil}
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}