package usecase

import (
	"back/dao"
	"back/model"
	"firebase.google.com/go/auth"
	"log"
)

func GetTweets(tags []string, id string) ([]model.Tweet, error) {
	tweets, err := dao.GetTweets(tags, id)

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

func PostReply(token *auth.Token, body string, tweetID string) (model.Tweet, error) {
	tweet, err := dao.PostReply(token, body, tweetID)

	if err != nil {
		log.Println("an error occurred at usecase/tweets")
	}
	return tweet, err
}

func GetUserTweets(userID string, tags []string) ([]model.Tweet, error) {
	tweets, err := dao.GetUserTweets(userID, tags)

	if err != nil {
		log.Println("an error occurred at usecase/tweets")
	}
	return tweets, err
}

func GetFollowingUserTweets(token *auth.Token, tags []string) ([]model.Tweet, error) {
	tweets, err := dao.GetFollowingUserTweets(token, tags)

	if err != nil {
		log.Println("an error occurred at usecase/tweets")
	}
	return tweets, err
}
