package model
import (
	"github.com/dgrijalva/jwt-go"
)
type Token struct {
	Name string
	Email string
	StandardClaims jwt.StandardClaims
}
// Need to validate this later
func (t Token) Valid() error {
	return nil
}

