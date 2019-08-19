package controllers

import (
	"auth-service-template/models"
	u "auth-service-template/utils"
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
	u.Respond(w, nil)
}
