package models

import (
	u "auth-service-template/utils"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type Contact struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	UserId string `json:"user_id"` //The user that this contact belongs to
}

/*
 This struct function validate the required parameters sent through the http request body

returns message and true if the requirement is met
*/
func (contact *Contact) Validate() (map[string]interface{}, bool) {

	if contact.Name == "" {
		return u.Message(false, "Contact name should be on the payload"), false
	}

	if contact.Phone == "" {
		return u.Message(false, "Phone number should be on the payload"), false
	}

	if contact.UserId == "" {
		return u.Message(false, "User is not recognized"), false
	}

	//All the required parameters are present
	return u.Message(true, "success"), true
}

func (contact *Contact) Create() map[string]interface{} {

	if resp, ok := contact.Validate(); !ok {
		return resp
	}

	db := GetDB()
	collection := db.Collection("contacts")
	_, err := collection.InsertOne(context.Background(), contact)

	if err == nil {
		resp := u.Message(true, "success")
		resp["contact"] = contact
		return resp
	} else {
		return nil
	}

}

func GetContact(id string) *Contact {

	db := GetDB()
	collection := db.Collection("contacts")
	contact := Contact{}
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(contact)

	if err != nil {
		return nil
	}
	return &contact
}

func GetContacts(user string) []*Contact {

	db := GetDB()
	collection := db.Collection("contacts")
	contacts := []*Contact{}
	cur, err := collection.Find(context.Background(), bson.M{"user_id": user})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for cur.Next(context.Background()) {

		// create a value into which the single document can be decoded
		var contact Contact
		err := cur.Decode(&contact)
		if err != nil {
			log.Fatal(err)
		}
		contacts = append(contacts, &contact)
	}

	return contacts
}
