package controller

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
var followDBService dbservice.FollowDBservice

func init() {
	log.Println("Loading .env if exists")
	godotenv.Load()
	log.Println("Initializing db connection")
	username := os.Getenv("PG_USERNAME")
	password := os.Getenv("PG_PASSWORD")
	databaseName := os.Getenv("PG_DBNAME")
	databaseHost := os.Getenv("DB_HOST")
	dbConnection := postgres.ConnectDB(username, password, databaseName, databaseHost)
	followDBService = dbservice.FollowDBservice{dbConnection}
}


// Create follow godoc
// @Summary Create a follow relationship
// @Description Create a follow relationship of two users
// @Accept  json
// @Produce  json
// @Param follow body model.Follow true "JSON describing relationship"
// @Success 201 {object} model.Follow
// @Failure 409 {object} utils.HttpError
// @Router /follow [POST]
func CreateFollow(w http.ResponseWriter, r *http.Request) {
	// Decode json body into model.Follow
	follow := &model.Follow{}
	json.NewDecoder(r.Body).Decode(follow)
	savedFollow, err := followDBService.CreateRelation(*follow)
	if err != nil {
		log.Println(err.Message)
		httpErr := utils.HttpError{
			Err: err.Message,
		}
		w.WriteHeader(err.Code)
		json.NewEncoder(w).Encode(httpErr)
		return

	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(savedFollow)
		return
	}

}

// Find the followers of the user with username
// @Summary Get an array of usernames that are following the user given in the parameter
// @Description Fetch all usernames following the username
// @Produce json
// @Param uesrname path string true "The user name to get the followers"
// @Success 200 {array} string
// @Failure 404 {object} utils.HttpError
// @Router /followers [GET]
func GetFollowers(w http.ResponseWriter, r* http.Request) {
	username, ok := r.URL.Query()["username"]
	if !ok {
		log.Println("Cannot get username from query ")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.HttpError{Err: "Failed to get parameter: username"})
		return
	}
	followers, err := followDBService.FindFollowers(username[0])
	if err != nil {
		log.Println(err)
		httpError := utils.HttpError{Err: err.Message}
		w.WriteHeader(err.Code)
		json.NewEncoder(w).Encode(httpError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(followers)

}

// Find the followees of the user with username
// @Summary Get an array of usernames that the user is following given in the parameter
// @Description Fetch all usernames the user is following
// @Produce json
// @Param uesrname path string true "The user name to get the followees"
// @Success 200 {array} string
// @Failure 404 {object} utils.HttpError
// @Router /followees [GET]
func GetFollowees(w http.ResponseWriter, r* http.Request) {
	username, ok := r.URL.Query()["username"]
	if !ok {
		log.Println("Cannot get username from query ")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.HttpError{Err: "Failed to get parameter: username"})
		return
	}
	followers, err := followDBService.FindFollowees(username[0])
	if err != nil {
		log.Println(err)
		httpError := utils.HttpError{Err: err.Message}
		w.WriteHeader(err.Code)
		json.NewEncoder(w).Encode(httpError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(followers)

}


