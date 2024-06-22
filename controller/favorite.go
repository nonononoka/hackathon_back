package controller

import (
	"back/usecase"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostFavorites(ctx *gin.Context) {
	tweetID := ctx.Query("id")
	token := ctx.MustGet("token").(*auth.Token)

	error := usecase.PostFavorites(token, tweetID)
	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

// ユーザーがlikeしてるツイート全部持ってくる
func GetFavorites(ctx *gin.Context) {
	token := ctx.MustGet("token").(*auth.Token)

	tweets, error := usecase.GetFavorites(token)
	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, tweets)
}

func DeleteFavorites(ctx *gin.Context) {
	tweetID := ctx.Query("id")
	token := ctx.MustGet("token").(*auth.Token)

	error := usecase.DeleteFavorites(token, tweetID)
	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
