package controller

import (
	"back/usecase"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostLike(ctx *gin.Context) {
	tweetID := ctx.Param("id")
	token := ctx.MustGet("token").(*auth.Token)

	error := usecase.PostLike(token, tweetID)
	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
