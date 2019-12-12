package models

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

/*
JWT claims struct
*/
type Token struct {
	UserId   string `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Username string `bson:"username,omitempty" json:"username,omitempty"`
	jwt.StandardClaims
	Role                   []string                 `bson:"role,omitempty" json:"role,omitempty"`
	OrganizationPermission []OrganizationPermission `bson:"organization_permission,omitempty" json:"organization_permission,omitempty"`
}


func CreateTokens(user User) (token *jwt.Token, refreshToken *jwt.Token) {
	tk := &Token{UserId: user.ID, Username: user.Username, Role: user.Role, OrganizationPermission: user.OrganizationPermission}
	tk.ExpiresAt = time.Now().Add(time.Minute * 24).Unix()

	token = jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tkRefresh := &Token{UserId: user.ID, Username:user.Username}
	tkRefresh.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()

	refreshToken = jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tkRefresh)
	return token, refreshToken
}

func CreateStringToken(token *jwt.Token) (string, error) {
	return token.SignedString([]byte(os.Getenv("token_password")))

}