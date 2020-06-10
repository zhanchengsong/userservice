package utils

import (
	"github.com/zhanchengsong/userservice/model"
	"golang.org/x/crypto/bcrypt"
	"time"
)
// TokenizeUser turns a user into JWTToken
func TokenizeUser(user model.User) string {
	// Setup expires time
	expiresAt := time.Now().Add(time.Minute * 100000).Unix()
	tk:= model.Token{}
}
