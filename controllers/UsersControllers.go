package controllers

import (
	u "auth-service-template/utils"
	"fmt"
	"net/http"

	"github.com/gorilla/context"
)

var ListUsers = func(w http.ResponseWriter, r *http.Request) {
	limit := r.FormValue("limit")
	skip := r.FormValue("skip")
	val := context.Get(r, "user")

	fmt.Println(limit)
	fmt.Println(skip)
	fmt.Println(val)
	u.Respond(w, nil)
}
