package controllers

import (
	u "auth-service-template/utils"
	"fmt"
	"net/http"
)

var ListUsers = func(w http.ResponseWriter, r *http.Request) {
	limit := r.FormValue("limit")
	skip := r.FormValue("skip")
	fmt.Println(limit)
	fmt.Println(skip)
	u.Respond(w, nil)
}
