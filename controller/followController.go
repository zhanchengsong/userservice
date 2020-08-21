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
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(httpErr)
		return

	} else {
		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(savedFollow)
		return
	}

}

