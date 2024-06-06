package usecase

import (
	"back/dao"
	"back/model"
	"firebase.google.com/go/auth"
	"log"
)

func GetMyTweets(token *auth.Token) ([]model.Tweet, error) {
	tweets, err := dao.GetMyTweets(token)

	if err != nil {
		log.Println("an error occurred at usecase/tweets")
	}
	return tweets, err
}

func PostTweet(token *auth.Token, body string) (model.Tweet, error) {
	tweet, err := dao.PostTweet(token, body)

	if err != nil {
		log.Println("an error occurred at usecase/tweets")
	}
	return tweet, err
}
