package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/zhanchengsong/userservice/model"
	"os"
	"testing"
)

func TestEncodeEmptyToJWT(t *testing.T) {
	t.Log("Test tokenize an empty user")
	t.Log("Setting JWT_SECRET to testsecret")
	os.Setenv("JWT_SECRET", "testsecret")
	emptyUser := model.User{}
	t.Log("Tokenizing User")
	jwtToken, err := TokenizeUser(emptyUser)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log("Parsing jwt string into token ")
	// Verify that jwtToken can be parsed
	token, perr := decodeJWT(jwtToken)
	if perr != nil {
		t.Error(perr.Error())
		return
	}
	// Try cast the token Claims into jwt MapClaims (concrete)
	// Type assertions
	t.Log("Parsing claims in jwt")
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Error("Failed to cast token into jwt.MapClaims")
	}
	t.Log("Checking fields")
	Email, ok := claims["Email"]
	if !ok {
		t.Error("Missing Email field")
	}
	if Email != "" {
		t.Error("Email field should not have a value ")
	}
	Name, ok := claims["Name"]
	if Name != "" {
		t.Error("Missing User Name field")
	}
	standardClaims, ok := claims["StandardClaims"]
	if !ok {
		t.Error("Error getting standard claims")
	}
	if standardClaims == nil {
		t.Error("Not able to get standard claims")
	}
	t.Log("Test create jwt for empty user successful")
}

func decodeJWT(jwtToken string) (jwt.Token, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("testsecret"), nil
	},
	)
	return *token, err
}
