package controllers

import (
	"auth-service-template/models"
	"encoding/json"
	"net/http"
)

var CreateUser = func(w http.ResponseWriter, r *http.Request) {

	account := &models.User{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		response := map[string]interface{}{"status": false, "message": "Invalid request"}
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	err = account.Create() //Create account
	var response models.UserResponse
	if err==nil {
		response = models.UserResponse{Code: 200, Message: "User saved", Data: []models.User{*account}}
	} else {
		response = models.UserResponse{Code: 400, Message: "error", Data: nil}
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &models.User{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := models.UserResponse{Code: 401, Message: err.Error(), Data: nil}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	var response models.UserResponse
	err, user := models.Login(account.Username, account.Password)

	if err  != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response = models.UserResponse{Code: 401, Message: err.Error(), Data: nil}
	}
	response = models.UserResponse{Code: 200, Message: "retrieved user", Data: []models.User{user}}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
