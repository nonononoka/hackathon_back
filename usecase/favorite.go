package usecase

import (
	"back/dao"
	"back/model"
	"firebase.google.com/go/auth"
	"log"
)

func PostFavorites(token *auth.Token, tweetID string) error {
	err := dao.PostFavorites(token, tweetID)

	if err != nil {
		log.Println("an error occurred at usecase/likes")
	}
	return err
}

func GetFavorites(token *auth.Token) ([]model.Tweet, error) {
	tweets, err := dao.GetFavorites(token)

	if err != nil {
		log.Println("an error occurred at usecase/likes")
	}
	return tweets, err
}

func DeleteFavorites(token *auth.Token, tweetID string) error {
	err := dao.DeleteFavorites(token, tweetID)

	if err != nil {
		log.Println("an error occurred at usecase/likes")
	}
	return err
}
