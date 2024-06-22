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

func UnfollowUser(token *auth.Token, userID string) error {
	err := dao.UnfollowUser(token, userID)

	if err != nil {
		log.Println(err.Error())
	}
	return err
}

func GetFollowingUsers(token *auth.Token) ([]model.User, error) {
	followingUsers, err := dao.GetFollowingUsers(token)

	if err != nil {
		log.Println(err.Error())
	}
	return followingUsers, err
}

func GetFollowedUsers(token *auth.Token) ([]model.User, error) {
	followedUsers, err := dao.GetFollowedUsers(token)

	if err != nil {
		log.Println(err.Error())
	}
	return followedUsers, err
}
