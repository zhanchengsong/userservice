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
func (dbservice UserDbservice) saveUser(user model.User) (model.User, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	} else {
		user.Password = string(pass)
	}
	dbservice.DbConnection.Create(user)
	return user, nil
}

func (dbservice UserDbservice) findUser(email, password string ) (model.User, error) {
	user := &model.User{}
	err := dbservice.DbConnection.Where("Email = ?", email).First(user).Error
	if err != nil {
		log.Fatal(err.Error())
		return *user, err
	}

}

