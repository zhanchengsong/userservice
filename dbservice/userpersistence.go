package dbservice

import (
	"github.com/jinzhu/gorm"
	"github.com/zhanchengsong/userservice/model"
	"golang.org/x/crypto/bcrypt"
	"log"
)
type UserDbservice struct {
	DbConnection *gorm.DB
}

// saveUser method saves a user entry in the database
func (dbservice *UserDbservice) SaveUser(user model.User) (model.User, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	} else {
		user.Password = string(pass)
	}
	saveerr := dbservice.DbConnection.Create(&user).Error
	if saveerr != nil {
		log.Println(saveerr)
		return user, saveerr
	}
	return user, nil
}

func (dbservice *UserDbservice) FindUser(email, password string ) (model.User, error) {
	user := &model.User{}
	err := dbservice.DbConnection.Where("Email = ?", email).First(user).Error
	if err != nil {
		log.Println(err.Error())
		return *user, err
	}
	// Verify password
	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf!=nil {
		log.Println(err.Error())
		return model.User{}, errf
	}
	return *user, nil
}

