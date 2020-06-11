package controllers

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/zhanchengsong/userservice/dbservice"
	"github.com/zhanchengsong/userservice/model"
	"github.com/zhanchengsong/userservice/postgres"
	"log"
	"net/http"
	"os"
)

type ErrorResponse struct {
	Err string
}

type error interface {
	Error() string
}

var userDBService dbservice.UserDbservice

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Cannot not load env file")
	}
	username := os.Getenv("PG_USERNAME")
	password := os.Getenv("PG_PASSWORD")
	databaseName := os.Getenv("PG_DBNAME")
	databaseHost := os.Getenv("DB_HOST")
	dbConnection := postgres.ConnectDB(username, password, databaseName, databaseHost)
	userDBService = dbservice.UserDbservice{dbConnection}
}
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Decode json body into model.User
	user := &model.User{}
	json.NewDecoder(r.Body).Decode(user)
	savedUser, err := userDBService.SaveUser(*user)
	if err != nil {
		log.Println(err.Error())
		httpErr:= ErrorResponse{
			Err: "Failed User registration",
		}
		encodingError := json.NewEncoder(w).Encode(httpErr)
		if encodingError != nil {
			log.Println(encodingError)
		}
	}
	savedUser.Password = ""
	json.NewEncoder(w).Encode(savedUser)
}
