package models

import (
	u "auth-service-template/utils"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	options "go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

/*
JWT claims struct
*/
type Token struct {
	UserId string
	jwt.StandardClaims
}

//a struct to rep user user
type User struct {
	Username string   `bson:"username" json:"username,omitempty"`
	Password string   `bson:"password" json:"password,omitempty"`
	Token    string   `bson:"token" json:"token,omitempty"`
	ID       string   `bson:"_id,omitempty" json:"_id,omitempty"`
	Email    string   `json: "email"`
	Role     []string `json: "role"`
}

//Validate incoming user details...
func (user *User) Validate() (map[string]interface{}, bool) {

	if user.Username == "" {
		return u.Message(false, "Username address is required"), false
	}

	if len(user.Password) < 2 {
		return u.Message(false, "Password is required"), false
	}

	//check for errors and duplicate Usernames
	db := GetDB()
	collection := db.Collection("users")
	foundUser := &User{}
	err := collection.FindOne(context.Background(), bson.M{"username": user.Username}).Decode(foundUser)

	if err != nil {
		fmt.Println(err)
	}

	if foundUser.Username != "" {
		return u.Message(false, "Username address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (user *User) Create() map[string]interface{} {

	if resp, ok := user.Validate(); !ok {
		return resp
	}
	fmt.Println(user)
	user.ID = bson.NewObjectId().Hex()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	db := GetDB()
	collection := db.Collection("users")
	res, err := collection.InsertOne(context.Background(), user)
	fmt.Println(res)

	if err != nil {
		return u.Message(false, "Failed to create user, connection error.")
	}

	//Create new JWT token for the newly registered user
	fmt.Println(user.ID)

	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	user.Password = "" //delete password

	response := u.Message(true, "User has been created")
	response["user"] = user
	return response
}

func Login(username, password string) (resp map[string]interface{}, code int) {
	user := &User{}
	db := GetDB()
	collection := db.Collection("users")
	foundUser := &User{}
	err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(foundUser)
	if err != nil {
		fmt.Println(err)
		return u.Message(false, "Invalid login credentials. Please try again"), 401
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again"), 401
	}
	//Worked! Logged In
	user.Password = ""
	user.Username = foundUser.Username
	//Create JWT token
	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString //Store the token in the response

	resps := u.Message(true, "Logged In")
	resps["user"] = user
	return resps, 200
}

func GetUser(u string) *User {

	acc := &User{}
	db := GetDB()
	collection := db.Collection("users")
	foundUser := User{}
	err := collection.FindOne(context.Background(), bson.M{"_d": foundUser.ID}).Decode(foundUser)
	if err != nil { //User not found!
		return nil
	}

	acc.Password = ""
	return acc
}

func GetUsers(limit int, skip int) []User {

	db := GetDB()
	collection := db.Collection("users")
	users := []User{}
	limitOption := int64(limit)
	skipOption := int64(skip)
	findOptions := options.FindOptions{Limit: &limitOption, Skip: &skipOption}
	cur, err := collection.Find(context.Background(), bson.M{}, &findOptions)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for cur.Next(context.Background()) {

		// create a value into which the single document can be decoded
		var user User
		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	return users
}
