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

func GetMeFollowingUsers(token *auth.Token) ([]model.User, error) {
	followingUsers, err := dao.GetMeFollowingUsers(token)

	if err != nil {
		log.Println(err.Error())
	}
	return followingUsers, err
}

func GetMeFollowedUsers(token *auth.Token) ([]model.User, error) {
	followedUsers, err := dao.GetMeFollowedUsers(token)

	if err != nil {
		log.Println(err.Error())
	}
	return followedUsers, err
}

func GetFollowingUsers(token *auth.Token, userID string) ([]model.User, error) {
	followingUsers, err := dao.GetFollowingUsers(token, userID)

	if err != nil {
		log.Println(err.Error())
	}
	return followingUsers, err
}

func GetFollowedUsers(token *auth.Token, userID string) ([]model.User, error) {
	followedUsers, err := dao.GetFollowedUsers(token, userID)

	if err != nil {
		log.Println(err.Error())
	}
	return followedUsers, err
}
