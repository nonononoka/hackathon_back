package controller

import (
	"back/usecase"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func FollowUser(ctx *gin.Context) {
	userID := ctx.Query("id")
	token := ctx.MustGet("token").(*auth.Token)
	followingUser, error := usecase.FollowUser(token, userID)

	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, followingUser)
}

func UnfollowUser(ctx *gin.Context) {
	userID := ctx.Query("id")
	token := ctx.MustGet("token").(*auth.Token)
	error := usecase.UnfollowUser(token, userID)

	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}

func GetFollowingUser(ctx *gin.Context) {
	token := ctx.MustGet("token").(*auth.Token)
	followingUsers, error := usecase.GetFollowingUsers(token)

	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, followingUsers)
}

func GetFollowedUser(ctx *gin.Context) {
	token := ctx.MustGet("token").(*auth.Token)
	followedUsers, error := usecase.GetFollowedUsers(token)

	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, followedUsers)
}
