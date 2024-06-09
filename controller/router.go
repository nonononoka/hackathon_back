package controller

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

func StartServer() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"DELETE",
			"PUT",
		},
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
		},
	}))

	v1 := router.Group("twitter/api/v1")
	v1.Use(authMiddleware())
	{
		v1.GET("me", GetMe)
		v1.POST("me", PostMe)
		v1.GET("tweets", GetTweets)
		v1.POST("tweet", PostTweet)
		v1.GET("tags", GetTags)
		// 特定のuserのツイートを取得
		v1.GET("users/:id/tweets", GetUserTweets)
		// 全部のuserを取得
		v1.GET("users", GetUsers)
		// :idのuserをフォロー follow/:idで良さそう
		v1.POST("friendships/:id/follow", FollowUser)
		// フォローしてるuserを取得 follow-usersで良さそう
		v1.GET("friendships/follow", GetFollowingUser)
		// フォローしてるuserのツイート一覧 follow-users/tweetsで良さそう
		v1.GET("friendships/tweets", GetFollowingUserTweets)
		// 特定ツイートのいいね
		//v1.POST("likes/:id/", PostLikes)
		//// 特定userがいいねしてるツイートを取得する
		//v1.GET("users/:id/likes-tweets", GetLikeTweets)
	}
	router.Run(":8080")
}
