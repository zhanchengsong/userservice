package dbservice

import (
	"github.com/jinzhu/gorm"
	myerror "github.com/zhanchengsong/userservice/error"
	"github.com/zhanchengsong/userservice/model"
	"log"
)

type FollowDBservice struct {
	DbConnection *gorm.DB
}

func (followDBService *FollowDBservice) CreateRelation(follow model.Follow) (model.Follow, *myerror.DBError) {
	// Get check

	var dupFollows []model.Follow
	getErr := followDBService.DbConnection.Where("Follower = ? AND Followee = ?", follow.Follower, follow.Followee).Find(&dupFollows).Error
	if getErr != nil {
		log.Println(getErr)
		return follow, &myerror.DBError{Code: 500, Message: getErr.Error(), Clause: "Failed to check existing relationship"}
	}
	if len(dupFollows) > 0 {
		return follow, &myerror.DBError{Code: 405, Message: "Duplicated entry", Clause: "Relationship exists"}
	}
	createErr := followDBService.DbConnection.Create(&follow).Error
	if createErr != nil {
		log.Println(createErr)
		return follow, &myerror.DBError{Code:500, Message: createErr.Error(), Clause: "Following failed"}
	}
	return follow,nil
}

//func (followDBservice *FollowDBservice) findFollowers(followee string) ([]string, *myerror.DBError) {
//
//}
