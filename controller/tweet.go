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
	tags := ctx.QueryArray("tags")
	tweets, error := usecase.GetTweets(tags)
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
	log.Printf("post tweet")
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
