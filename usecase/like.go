package usecase

import (
	"back/dao"
	"back/model"
	"firebase.google.com/go/auth"
	"log"
)

func PostLike(token *auth.Token, tweetID string) error {
	err := dao.PostLike(token, tweetID)

	if err != nil {
		log.Println("an error occurred at usecase/likes")
	}
	return err
}

func GetLike(token *auth.Token) ([]model.Tweet, error) {
	tweets, err := dao.GetLike(token)

	if err != nil {
		log.Println("an error occurred at usecase/likes")
	}
	return tweets, err
}

func DeleteLike(token *auth.Token, tweetID string) error {
	err := dao.DeleteLike(token, tweetID)

	if err != nil {
		log.Println("an error occurred at usecase/likes")
	}
	return err
}
