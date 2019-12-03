package main

import (
	"auth-service-template/controllers"
	"auth-service-template/middlewares"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const logPath = "logs.log"

var Logger *log.Logger

func main() {

	godotenv.Load()

	InitLogger()

	router := mux.NewRouter()

	router.HandleFunc("/api/user/register", controllers.CreateUser).Methods("POST")

	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")

	router.HandleFunc("/api/user/list", controllers.ListUsers).Queries("limit", "{limit}", "skip", "{skip}").Methods("GET")

	router.Use(middlewares.MiddlewareLogger)  //attach JWT auth middleware
	router.Use(middlewares.JwtAuthentication) //attach JWT auth middleware

	//router.NotFoundHandler = app.NotFoundHandler

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	log.Println("Starting server on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))

}

func InitLogger() {

	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalln("Failed to open log file")
	}
	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)

}
