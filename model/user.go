package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model           // This is bascially an "extend"
	DisplayName string   `json:"displayName,omitempty"`
	Username    string   `json:"username",gorm:"type:varchar(100);unique_index"`
	Email       string   `gorm:"type:varchar(100);unique_index"`
	Password    string   `json:"-"`
	IconUrl    	string   `json:"iconUrl,omitempty"`
}
