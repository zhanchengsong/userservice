package controllers

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/zhanchengsong/userservice/dbservice"
	"github.com/zhanchengsong/userservice/model"
	"github.com/zhanchengsong/userservice/postgres"
	"github.com/zhanchengsong/userservice/utils"
	"log"
	"net/http"
	"os"
)

type ErrorResponse struct {
	Err string
}

type error interface {
	Error() string
	Cause() string
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
		log.Println(err.Message)
		httpErr:= ErrorResponse{
			Err: "Failed User registration",
		}
		w.WriteHeader(err.Code)
		json.NewEncoder(w).Encode(httpErr)
		return

	} else {
		savedUser.Password = ""
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(savedUser)
		return
	}

}

func Login(w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	json.NewDecoder(r.Body).Decode(&user)
	userFound, err := userDBService.FindUser(user.Email, user.Password)
	if err != nil {
		log.Println(err.Message)
		httpErr := ErrorResponse{
			Err: err.Message,
		}
		w.WriteHeader(err.Code)
		json.NewEncoder(w).Encode(httpErr)
		return
	}
	// Encode JWT if user is found
	jwt, errt := utils.TokenizeUser(userFound)
	if errt != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			struct{
				Err string
			}{errt.Error()},
			)
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(
			struct {
				JWTToken string
			}{
				jwt,
			})
	}
}

func findUserById(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.URL.Query()["userId"]
	if !ok {
		log.Println("userId is not found in the query parameter")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			struct{
				Err string
			}{
				Err: "Cannot get userId from query parameter",
			})
	}
	// Find the user

}