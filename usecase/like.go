package usecase

import (
	"back/dao"
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
