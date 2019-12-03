package controllers

import (
	"auth-service-template/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
)

var ListUsers = func(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	skip, _ := strconv.Atoi(r.FormValue("skip"))
	user := context.Get(r, "user")

	users := models.GetUsers(limit, skip)
	fmt.Println(limit)
	fmt.Println(skip)
	fmt.Println(users)
	fmt.Println(user)
	response := models.UserResponse{Code: 200, Message: "users retrieved", Data: users}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
