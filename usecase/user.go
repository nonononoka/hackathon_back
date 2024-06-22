package usecase

import (
	"back/dao"
	"back/model"
	"firebase.google.com/go/auth"
	"log"
)

func GetMe(token *auth.Token) (model.User, error) {
	userInfo, err := dao.GetMe(token)

	if err != nil {
		log.Println("an error occurred at usecase/user")
	}
	return userInfo, err
}

func PostMe(token *auth.Token) (model.User, error) {
	userInfo, err := dao.PostMe(token)

	if err != nil {
		log.Println("an error occurred at usecase/user", err)
	}
	return userInfo, err
}

func PutMe(token *auth.Token, name string, bio string, image string) (model.User, error) {
	userInfo, err := dao.PutMe(token, name, bio, image)

	if err != nil {
		log.Println("an error occurred at usecase/user", err)
	}
	return userInfo, err
}

func GetUsers(token *auth.Token) ([]model.User, error) {
	userInfos, err := dao.GetUsers(token)
	if err != nil {
		log.Println("an error occurred at usecase/user")
	}
	return userInfos, err
}
