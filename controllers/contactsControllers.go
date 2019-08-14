package controllers

import (
	"auth-service-template/models"
	u "auth-service-template/utils"
	"encoding/json"
	"net/http"
)

var CreateContact = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(string) //Grab the id of the user that send the request
	contact := &models.Contact{}

	err := json.NewDecoder(r.Body).Decode(contact)

	if err != nil {
		response := map[string]interface{}{"status": false, "message": "Error while decoding request body"}
		u.Respond(w, response)
		return
	}

	contact.UserId = user
	resp := contact.Create()
	u.Respond(w, resp)
}

var GetContactsFor = func(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value("user").(string)
	data := models.GetContacts(id)
	response := map[string]interface{}{"status": true, "message": "success"}
	response["data"] = data
	u.Respond(w, response)
}
