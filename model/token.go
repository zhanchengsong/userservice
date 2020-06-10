package model
import (
	"github.com/dgrijalva/jwt-go"
)
type Token struct {
	userID string
	Name string
	Email string
	StandardClaims *jwt.StandardClaims
}
