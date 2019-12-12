package controllers

import (
	"auth-service-template/models"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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

func Refresh(w http.ResponseWriter, r *http.Request) {
	log.Println("refreshing token")

	tokenHeader := r.Header.Get("Authorization") //Grab the token from the header
	var response models.RefreshResponse

	if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
		log.Println("error validating token. No token")
		response = models.RefreshResponse{Code: 401, Message: "missing auth token"}
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
	if len(splitted) != 2 {
		log.Println("error validating token, malformed token.")
		response = models.RefreshResponse{Code: 401, Message: "malformed token"}
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	tokenPart := splitted[1] //Grab the token part, what we are truly interested in

	claims := &models.Token{}
	tkn, err :=  jwt.ParseWithClaims(tokenPart, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("error validating token. Error in signature. "+ err.Error())

			w.WriteHeader(http.StatusUnauthorized)
			response = models.RefreshResponse{Code: 401, Message: "error in token "+ err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}
		log.Println("error validating token. "+ err.Error())
		w.WriteHeader(http.StatusBadRequest)
		response = models.RefreshResponse{Code: 401, Message: "missing auth token "+ err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("error validating token. Token invalid")
		response = models.RefreshResponse{Code: 401, Message: "token invalid "}
		json.NewEncoder(w).Encode(response)
		return
	}
	// (END) The code up-till this point is the same as the first part of the `Welcome` route

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Hour {
		w.WriteHeader(http.StatusBadRequest)
		response = models.RefreshResponse{Code: 401, Message: "too soon to refresh token"}
		json.NewEncoder(w).Encode(response)
		return
	}
	log.Println("creating token")

	user := models.GetUser(claims.Username)

	newToken, newRefreshToken := models.CreateTokens(*user)
	tokenString, err := models.CreateStringToken(newToken)
	refreshTokenString , err2 := models.CreateStringToken(newRefreshToken)

	if err != nil || err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response = models.RefreshResponse{Code: 401, Message: "missing auth token "+ err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}


	response = models.RefreshResponse{Code: 200, Message: "retrieved token", Token:tokenString, RefreshToken:refreshTokenString}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}