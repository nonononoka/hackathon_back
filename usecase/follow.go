package usecase

import (
	"back/dao"
	"back/model"
	"firebase.google.com/go/auth"
	"log"
)

func FollowUser(token *auth.Token, userID string) (model.User, error) {
	followingUser, err := dao.FollowUser(token, userID)

	if err != nil {
		log.Println(err.Error())
	}
	return followingUser, err
}

func GetFollowingUsers(token *auth.Token) ([]model.User, error) {
	followwingUsers, err := dao.GetFollowingUsers(token)

	if err != nil {
		log.Println(err.Error())
	}
	return followwingUsers, err
}
