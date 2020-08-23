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

func (followDBservice *FollowDBservice) CreateRelation(follow model.Follow) (model.Follow, *myerror.DBError) {
	// Get check

	var dupFollows []model.Follow
	getErr := followDBservice.DbConnection.Where("Follower = ? AND Followee = ?", follow.Follower, follow.Followee).Find(&dupFollows).Error
	if getErr != nil {
		log.Println(getErr)
		return follow, &myerror.DBError{Code: 500, Message: getErr.Error(), Clause: "Failed to check existing relationship"}
	}
	if len(dupFollows) > 0 {
		return follow, &myerror.DBError{Code: 409, Message: "Duplicated entry", Clause: "Relationship exists"}
	}
	createErr := followDBservice.DbConnection.Create(&follow).Error
	if createErr != nil {
		log.Println(createErr)
		return follow, &myerror.DBError{Code:500, Message: createErr.Error(), Clause: "Following failed"}
	}
	return follow,nil
}

func (followDBservice *FollowDBservice) FindFollowers(username string) ([]string, *myerror.DBError) {
	var followers []model.Follow
	getErr := followDBservice.DbConnection.Where("Followee = ?", username).Find(&followers).Error
	if getErr != nil {
		log.Println(getErr)
		return []string{}, &myerror.DBError{Code:500, Message: getErr.Error(), Clause: "Failed to get followers"}
	}
	var result []string
	for i := 0 ; i < len(followers) ; i ++ {
		result = append(result,followers[i].Follower)
	}
	return result, nil
}

func (followDBservice *FollowDBservice) FindFollowees(username string) ([]string, *myerror.DBError) {
	var followees []model.Follow
	getErr := followDBservice.DbConnection.Where("Follower = ?", username).Find(&followees).Error
	if getErr != nil {
		log.Println(getErr)
		return []string{}, &myerror.DBError{Code:500, Message: getErr.Error(), Clause: "Failed to get followees"}
	}
	var result []string
	for i := 0 ; i < len(followees) ; i ++ {
		result = append(result,followees[i].Followee)
	}
	return result, nil
}
