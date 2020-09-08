package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/zhanchengsong/userservice/model"
	"os"
	"time"
)

// TokenizeUser turns a user into JWTToken

func TokenizeUser(user model.User) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	// Setup expires time
	expiresAt := time.Now().Add(time.Minute * 100000).Unix()
	tk := model.Token{
		Name:     user.DisplayName,
		Email:    user.Email,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
			Issuer:    "twitterUserService",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenStr, nil

}
