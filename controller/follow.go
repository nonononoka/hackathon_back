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

func GetMeFollowingUsers(ctx *gin.Context) {
	token := ctx.MustGet("token").(*auth.Token)
	followingUsers, error := usecase.GetMeFollowingUsers(token)

	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, followingUsers)
}

func GetMeFollowedUsers(ctx *gin.Context) {
	token := ctx.MustGet("token").(*auth.Token)
	followedUsers, error := usecase.GetMeFollowedUsers(token)

	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, followedUsers)
}

func GetFollowedUsers(ctx *gin.Context) {
	userID := ctx.Param("id")
	token := ctx.MustGet("token").(*auth.Token)
	followedUsers, error := usecase.GetFollowedUsers(token, userID)

	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, followedUsers)
}

func GetFollowingUsers(ctx *gin.Context) {
	userID := ctx.Param("id")
	token := ctx.MustGet("token").(*auth.Token)
	followingUsers, error := usecase.GetFollowingUsers(token, userID)

	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, followingUsers)
}
