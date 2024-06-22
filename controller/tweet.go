package controller

import (
	"back/usecase"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Body struct {
	Text string   `json:"body"`
	Tags []string `json:"tags"`
}

func GetTweets(ctx *gin.Context) {
	token := ctx.MustGet("token").(*auth.Token)
	tags := ctx.QueryArray("tags")
	id := ctx.Query("id")
	tweets, error := usecase.GetTweets(token, tags, id)
	log.Println(tweets)

	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, tweets)
}

func PostTweet(ctx *gin.Context) {
	token := ctx.MustGet("token").(*auth.Token)

	var body Body
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}

	tweet, error := usecase.PostTweet(token, body.Text, body.Tags)
	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, tweet)
}

func PostReply(ctx *gin.Context) {
	token := ctx.MustGet("token").(*auth.Token)
	tweetID := ctx.Param("id") // replyするツイート
	var body Body
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}

	tweet, error := usecase.PostReply(token, body.Text, body.Tags, tweetID)
	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, tweet)
}

func GetUserTweets(ctx *gin.Context) {
	token := ctx.MustGet("token").(*auth.Token)
	tags := ctx.QueryArray("tags")
	userID := ctx.Param("id")
	tweets, error := usecase.GetUserTweets(token, userID, tags)
	log.Println(tweets)

	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, tweets)
}

func GetFollowingUserTweets(ctx *gin.Context) {
	tags := ctx.QueryArray("tags")
	token := ctx.MustGet("token").(*auth.Token)
	tweets, error := usecase.GetFollowingUserTweets(token, tags)

	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, tweets)
}

// // idのツイートの一連のツイートを返す。
func GetThreadTweets(ctx *gin.Context) {
	tweetID := ctx.Param("id")
	token := ctx.MustGet("token").(*auth.Token)
	threads, error := usecase.GetThreadTweets(token, tweetID)

	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, threads)
}
