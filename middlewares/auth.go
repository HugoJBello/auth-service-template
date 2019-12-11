package middlewares

import (
	"auth-service-template/models"
	"encoding/json"
	"github.com/gorilla/context"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{"/api/user/register", "/api/user/login"} //List of endpoints that doesn't require auth
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

			response = map[string]interface{}{"status": false, "message": "missing auth token"}
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			response = map[string]interface{}{"status": false, "message": "Invalid/Malformed auth token"}
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

		if err != nil { //Malformed token, returns with http code 403 as usual
			response = map[string]interface{}{"status": false, "message": "Invalid/Malformed auth token"}
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

		user := models.GetUser(tk.Username)

		if strings.Contains(requestPath, "/admin/") && !checkAdminRoute(requestPath, *user) {
			response = map[string]interface{}{"status": true, "message": "You are not Admin."}
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		setUserDataInContext(user, r)

		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}

func checkAdminRoute(requestPath string, user models.User) bool {
	return isAdminRole(user.Role)
}

func setUserDataInContext(user *models.User, r *http.Request){
	permissions := make(map[string]string)
	for _, organization := range user.OrganizationPermission{
		permissions[organization.OrganizationId] = organization.OrganziationRole
	}
	context.Set(r, "organization_permissions", permissions)
	context.Set(r, "username", user.Username)
	context.Set(r, "role", user.Role)
	context.Set(r, "email", user.Email)
	context.Set(r, "userId", user.ID)

}

func isAdminRole(role []string) bool {
	for _, b := range role {
		if "ADMIN" == b {
			return true
		}
	}
	return false
}
