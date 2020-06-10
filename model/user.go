package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model           // This is bascially an "extend"
	DisplayName string   `json:"displayName"`
	Username    string   `json:"username"`
	Email       string   `gorm:"type:varchar(100);unique_index"`
	Password    string   `json:"Password"`
	ID          string   `json:-`
}
