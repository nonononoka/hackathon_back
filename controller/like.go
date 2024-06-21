package controller

import (
	"back/usecase"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostLike(ctx *gin.Context) {
	tweetID := ctx.Query("id")
	token := ctx.MustGet("token").(*auth.Token)

	error := usecase.PostLike(token, tweetID)
	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

// ユーザーがlikeしてるツイート全部持ってくる
func GetLike(ctx *gin.Context) {
	token := ctx.MustGet("token").(*auth.Token)

	tweets, error := usecase.GetLike(token)
	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, tweets)
}

func DeleteLike(ctx *gin.Context) {
	tweetID := ctx.Query("id")
	token := ctx.MustGet("token").(*auth.Token)

	error := usecase.DeleteLike(token, tweetID)
	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
