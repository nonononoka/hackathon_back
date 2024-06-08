package usecase

import (
	"back/dao"
	"back/model"
	"firebase.google.com/go/auth"
	"log"
)

func GetTweets(token *auth.Token, tags []string) ([]model.Tweet, error) {
	tweets, err := dao.GetTweets(token, tags)

	if err != nil {
		log.Println("an error occurred at usecase/tweets")
	}
	return tweets, err
}

func PostTweet(token *auth.Token, body string, tags []string) (model.Tweet, error) {
	tweet, err := dao.PostTweet(token, body, tags)

	if err != nil {
		log.Println("an error occurred at usecase/tweets")
	}
	return tweet, err
}
