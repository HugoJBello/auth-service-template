package middlewares

import (
	"auth-service-template/models"
	"encoding/json"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"log"
	"net/http"
	"os"
	"strings"
)

var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{"/api/user/register","/api/user/refresh", "/api/user/login"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path                                    //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			log.Println("missing token")
			response = map[string]interface{}{"status": false, "message": "missing auth token"}
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			log.Println("invalid token")
			response = map[string]interface{}{"status": false, "message": "Invalid/Malformed auth token. Tokens must be Bearer"}
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil {
			log.Println(err)
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				response = map[string]interface{}{"status": false, "message": "Invalid/Old auth token "+ err.Error()}
				w.Header().Add("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}
			response = map[string]interface{}{"status": false, "message": "Invalid/Malformed auth token. "+ err.Error()}
			log.Println(err)
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			response = map[string]interface{}{"status": true, "message": "Token is not valid."}
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		//user := models.GetUser(tk.Username)
		err = setUserDataInContext(*tk, r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			response = map[string]interface{}{"status": false, "message": "Invalid/Old auth token "+ err.Error()}
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		if strings.Contains(requestPath, "/admin/") && !checkAdminRoute(requestPath, tk.Role) {
			response = map[string]interface{}{"status": true, "message": "You are not Admin."}
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}

func checkAdminRoute(requestPath string, roles []string) bool {
	return isAdminRole(roles)
}

func setUserDataInContext(token models.Token, r *http.Request) error{
	if (token.Role == nil || token.OrganizationPermission == nil || token.Username == "" || token.UserId == ""){
		return errors.New("invalid token, token lacks the required fields")
	}

	permissions := make(map[string]string)
	for _, organization := range token.OrganizationPermission{
		permissions[organization.OrganizationId] = organization.OrganziationRole
	}

	context.Set(r, "organization_permissions", permissions)
	context.Set(r, "username", token.Username)
	context.Set(r, "role", token.Role)
	context.Set(r, "userId", token.UserId)
	return nil
}

func isAdminRole(role []string) bool {
	for _, b := range role {
		if "ADMIN" == b {
			return true
		}
	}
	return false
}