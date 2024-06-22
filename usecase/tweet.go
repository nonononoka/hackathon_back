package usecase

import (
	"back/dao"
	"back/model"
	"firebase.google.com/go/auth"
	"log"
)

func GetTweets(token *auth.Token, tags []string, id string) ([]model.Tweet, error) {
	tweets, err := dao.GetTweets(token, tags, id)

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

func PostReply(token *auth.Token, body string, tags []string, tweetID string) (model.Tweet, error) {
	tweet, err := dao.PostReply(token, body, tags, tweetID)

	if err != nil {
		log.Println("an error occurred at usecase/tweets")
	}
	return tweet, err
}

func GetUserTweets(token *auth.Token, userID string, tags []string) ([]model.Tweet, error) {
	tweets, err := dao.GetUserTweets(token, userID, tags)

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

func GetThreadTweets(token *auth.Token, tweetID string) ([]model.Tweet, error) {
	threadTweets, err := dao.GetThreadTweets(token, tweetID)

	if err != nil {
		log.Println("an error occurred at usecase/tweets")
	}
	return threadTweets, err
}
