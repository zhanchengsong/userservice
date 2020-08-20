package model
import (
	"github.com/jinzhu/gorm"
)
// A simple structure mapping the following relationship
type Follow struct {
	gorm.Model           // This is bascially an "extend"
	Follower    string   `json:"follower"`
	Followee    string   `json:"followee"`
}
