package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	guuid "github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/zhanchengsong/userservice/aws-service"
	"github.com/zhanchengsong/userservice/dbservice"
	"github.com/zhanchengsong/userservice/model"
	"github.com/zhanchengsong/userservice/postgres"
	"github.com/zhanchengsong/userservice/utils"
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
	log.Println("Loading .env if exists")
	godotenv.Load()
	log.Println("Initializing db connection")
	username := os.Getenv("PG_USERNAME")
	password := os.Getenv("PG_PASSWORD")
	databaseName := os.Getenv("PG_DBNAME")
	databaseHost := os.Getenv("DB_HOST")
	dbConnection := postgres.ConnectDB(username, password, databaseName, databaseHost)
	userDBService = dbservice.UserDbservice{dbConnection}
}

// Create User With Icon Doc
// @Summary Create a user
// @Description Create a user in the database, this call accepts a multiform with image file as icon
// @Accept  mpfd
// @Produce  json
// @Param username formData string true "Username of the registering user, should contain spaces"
// @Param displayName formData string true "A user friendly will be displayed on the UI"
// @Param email formData string true "The email address of the user"
// @Param password formData string true "Password of the user"
// @Param icon formData file true "Icon picture uploaded will be used in user profile"
// @Success 201 {object} model.User
// @Failure 409 {object} utils.HttpError
// @Failure 500 {object} utils.HttpError
// @Router /userWithIcon [POST]
func CreateUserWithIcon(w http.ResponseWriter, r *http.Request) {
	// Upload user icons
	supportedExtensions := []string{"jpeg", "jpg", "png"}
	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("icon")
	if err != nil {
		err := utils.HttpError{Err: "Failed to read icon file"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}
	defer file.Close()
	names := strings.Split(header.Filename, ".")
	log.Println("Received Filename: ", names[0])
	extension := names[1]
	// If the extension is not supported
	if !(sort.SearchStrings(supportedExtensions, extension) < len(supportedExtensions)) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(utils.HttpError{Err: "Unsupported format for icon"})
	}
	filename := guuid.New().String() + "." + extension
	location, err := aws_service.UploadIconToS3(file, filename)
	if err != nil {
		log.Println(err)

		httpErr := utils.HttpError{
			Err: err.Error(),
		}
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(httpErr)
		return
	}
	log.Println(location)
	// Start dealing with other form data
	username := r.FormValue("username")
	displayName := r.FormValue("displayName")
	email := r.FormValue("email")
	password := r.FormValue("password")
	icon := filename
	// Create User object from these values
	userCreate := model.User{Username: username, DisplayName: displayName, Email: email, Password: password, IconUrl: icon}
	createdUser, saveerr := userDBService.SaveUser(userCreate)
	if saveerr != nil {
		log.Println(saveerr.Message)
		httpErr := utils.HttpError{
			Err: saveerr.Message,
		}
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(httpErr)
		return

	} else {
		createdUser.Password = ""
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdUser)
		return
	}
	return
}

// Create User Doc
// @Summary Create a user
// @Description Create a user in the database
// @Accept  json
// @Produce  json
// @Param user body model.User true "JSON body describing user"
// @Success 201 {object} model.User
// @Failure 409 {object} utils.HttpError
// @Failure 500 {object} utils.HttpError
// @Router /user [POST]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Decode json body into model.User
	user := &model.User{}
	json.NewDecoder(r.Body).Decode(user)
	savedUser, err := userDBService.SaveUser(*user)
	if err != nil {
		log.Println(err.Message)
		httpErr := utils.HttpError{
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

// Create User Token after login
// @Summary Create a JWToken for user login and return the entire profile
// @Description Generate a JWToken if username/password is stored in database and return the complete profile including JWT Token
// @Accept  json
// @Produce  json
// @Param user body model.User true "A body describing user profile including jwtToken"
// @Success 201 {object} model.User
// @Failure 409 {object} utils.HttpError
// @Failure 500 {object} utils.HttpError
// @Router /login [POST]
func Login(w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	json.NewDecoder(r.Body).Decode(&user)
	userFound, err := userDBService.FindUser(user.Email, user.Password)
	if err != nil {
		log.Println(err.Message)
		httpErr := utils.HttpError{
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
			struct {
				Err string
			}{errt.Error()},
		)
	} else {
		userFound.JWTToken = jwt
		userFound.Password = ""
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(userFound)
	}
}

// Find the user by either user Id or user name
// @Summary Create a JWToken for user login and return the entire profile
// @Description Generate a JWToken if username/password is stored in database and return the complete profile including JWT Token
// @Produce  json
// @Param userId path string false "User ID to look for"
// @Param username path string false "Username to look for"
// @Success 200 {object} model.User
// @Failure 404 {object} utils.HttpError
// @Router /user [GET]
func FindUser(w http.ResponseWriter, r *http.Request) {
	_, idOk := r.URL.Query()["userId"]
	if idOk {
		FindUserById(w, r)
	} else {
		_, nameOk := r.URL.Query()["username"]
		if nameOk {
			FindUserByUsername(w, r)
		}
	}
}
// Check the user by username
// @Summary Check if a user exists
// @Description check if a user exists by the username
// @Produce json
// @Param username path string false "Username to look for"
// @Failure 404 {object} utils.HttpError
// @Router /user [GET]
func CheckUserByName(w http.ResponseWriter, r *http.Request) {
	username, ok := r.URL.Query()["username"]
	if !ok {
		log.Println("Failed getting username param from url")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			struct {
				Err string
			}{
				Err: "Missing path parameter username",
			})
	}
	found, err := userDBService.CheckUserByUsername(username[0]);
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			struct {
				Err string
			}{
				Err: "Failed when search user",
			})
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct {
			Exists bool
		}{
			Exists: found,
		})
	}
}

func CheckUserByEmail(w http.ResponseWriter, r *http.Request) {
	email, ok := r.URL.Query()["email"]
	if !ok {
		log.Println("Failed getting email param from url")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			struct {
				Err string
			}{
				Err: "Missing path parameter email",
			})
	}
	found, err := userDBService.CheckUserByEmail(email[0])
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			struct {
				Err string
			}{
				Err: "Failed when search user",
			})
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct {
			Exists bool
		}{
			Exists: found,
		})
	}
}
func FindUserById(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.URL.Query()["userId"]
	if !ok {
		log.Println("userId is not found in the query parameter")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			struct {
				Err string
			}{
				Err: "Cannot get userId from query parameter",
			})
	}
	// Find the user
	userIdInt, _ := strconv.Atoi(userId[0])
	foundUser, err := userDBService.FindUserById(userIdInt)
	if err != nil {
		w.WriteHeader(err.Code)
		json.NewEncoder(w).Encode(err)
		return
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(foundUser)
	}
	return

}

// Find the usernames starting from certain prefix
// @Summary Get an array of usernames that start with prefix
// @Description Fetch all usernames starting with the input prefix
// @Produce json
// @Param uesrprefix path string true "Prefix in the username to search for"
// @Success 200 {array} string
// @Failure 404 {object} utils.HttpError
// @Router /users [GET]
func FindUsersByPrefix(w http.ResponseWriter, r *http.Request) {
	userprefix, ok := r.URL.Query()["userprefix"]
	if !ok {
		log.Println("userprefix is not found in the query parameter")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			struct {
				Err string
			}{
				Err: "Cannot get userprefix from query parameter",
			})
	}
	foundUsers, err := userDBService.FindUserByPrefix(userprefix[0])
	if err != nil {
		w.WriteHeader(err.Code)
		json.NewEncoder(w).Encode(err)
		return
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(foundUsers)
	}
	return

}

func FindUserByUsername(w http.ResponseWriter, r *http.Request) {
	username, ok := r.URL.Query()["username"]
	if !ok {
		log.Println("Failed getting username param from url")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(
			struct {
				Err string
			}{
				Err: "Cannot find the user with username",
			})
	}
	foundUser, err := userDBService.FindUserByUsername(username[0])
	if err != nil {
		log.Println(err.Message)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			struct {
				Err string
			}{
				Err: "Failed when search user",
			})
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(foundUser)
	}

}
