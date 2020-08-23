package dbservice

import (
	"fmt"
	"github.com/jinzhu/gorm"
	myerror "github.com/zhanchengsong/userservice/error"
	"github.com/zhanchengsong/userservice/model"
	"golang.org/x/crypto/bcrypt"
	"log"
)
type UserDbservice struct {
	DbConnection *gorm.DB
}

// saveUser method saves a user entry in the database
func (dbservice *UserDbservice) SaveUser(user model.User) (model.User, *myerror.DBError) {
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, &myerror.DBError{Code:500, Clause: "Cannot Hash", Message: err.Error()}
	} else {
		user.Password = string(pass)
	}
	saveerr := dbservice.DbConnection.Create(&user).Error
	if saveerr != nil {
		log.Println(saveerr)
		dbError := myerror.DBError{Code: 409, Message: saveerr.Error(), Clause: "DUPLICATE"}
		return user, &dbError
	}
	return user, nil
}

func (dbservice *UserDbservice) FindUser(email, password string ) (model.User, *myerror.DBError) {
	user := &model.User{}
	err := dbservice.DbConnection.Where("Email = ?", email).First(user).Error
	if err != nil {
		log.Println(err.Error())
		dbError := myerror.DBError{Clause: "Not Found", Code: 404, Message: err.Error()}
		return *user, &dbError
	}
	// Verify password
	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf!=nil {
		log.Println(errf.Error())
		dbError := myerror.DBError{Clause: "Password Error", Code: 401, Message: errf.Error()}
		return model.User{}, &dbError
	}
	return *user, nil
}

func (dbservice *UserDbservice) FindUserById(userId int) (model.User, *myerror.DBError) {
	user := &model.User{}
	err := dbservice.DbConnection.Where("ID = ?", userId).First(user).Error
	if err != nil {
		log.Println(err.Error())
		dbError := myerror.DBError{Clause: "User Not Found", Code: 404, Message: err.Error()}
		return model.User{}, &dbError
	}
	user.Password = ""
	return *user, nil
}
// FindUserByPrefix find the list of users containing the name of prefix
func (dbservice *UserDbservice) FindUserByPrefix(userNamePrefix string) ([]model.User, *myerror.DBError) {
	users := &( []model.User{} )
	err := dbservice.DbConnection.Where("username LIKE ?", fmt.Sprintf("%%%s%%", userNamePrefix)).Find(&users).Error
	if err != nil {
		log.Println(err.Error())
		dbError := myerror.DBError{Clause: "Users Not Found", Code: 404, Message: err.Error()}
		return *users, &dbError
	}
	// Clean out the data
	for _, user := range *users {
		user.Password = ""
	}
	return *users, nil
}

// Get user profile
func (dbservice *UserDbservice) FindUserByUsername(username string) (model.User, *myerror.DBError) {
	user := &( model.User{})
	err := dbservice.DbConnection.Where("username = ?", username).First(&user).Error
	if err != nil {
		log.Println(err.Error())
		dbError := myerror.DBError{Clause: "Cannot find user", Code: 404, Message: err.Error()}
		return *user, &dbError
	}
	return *user, nil
}