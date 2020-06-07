package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model           // This is bascially an "extend"
	DisplayName string   `json:"displayName"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Password    string   `json:"Password"`
	ID          string   `json:-`
	Friends     []string `json:-`
}
