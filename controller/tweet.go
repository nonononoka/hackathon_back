package controller

import (
	"back/usecase"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Content struct {
	Text string `json:"text"`
}

func GetMyTweets(ctx *gin.Context) {
	token := ctx.MustGet("token").(*auth.Token)

	tweets, error := usecase.GetMyTweets(token)
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
	var body Content
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}

	tweet, error := usecase.PostTweet(token, body.Text)
	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, tweet)
}
