package controller

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Todo struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

func StartServer() {
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "https://tech-tweet.vercel.app"}
	corsConfig.AddAllowHeaders("Authorization")
	router.Use(cors.New(corsConfig))

	v1 := router.Group("twitter/api/v1")
	v1.Use(authMiddleware())
	{
		v1.GET("hello", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"message": "hello"}) })
		v1.GET("me", GetMe)
		v1.POST("me", PostMe)
		v1.PUT("me", PutMe)
		// /tweets/tags=?みたいな
		v1.GET("tweets", GetTweets)
		// 一つのツイートのスレッドを返す。投稿された順番になるように。
		v1.GET("threads/:id", GetThreadTweets)
		v1.POST("tweet", PostTweet)
		v1.GET("tags", GetTags)
		// 特定のuserのツイートを取得 /users/:id/tweets
		v1.GET("users/:id/tweets", GetUserTweets)
		// 全部のuserを取得 /users
		v1.GET("users", GetUsers)
		// 特定ツイートのいいね /:idのツイートをファボする。
		v1.POST("tweets/favorites", PostFavorites)
		// これいらないかも
		v1.GET("tweets/favorites", GetFavorites)
		v1.DELETE("tweets/favorites", DeleteFavorites)
		//// 特定userがいいねしてるツイートを取得する /users/:userid/tweets/favorites
		//v1.GET("users/:id/likes-tweets", GetLikeTweets)

		// 特定ツイートにリプライする。
		v1.POST("tweets/:id/reply", PostReply)

		// FF関係
		// :idのuserをフォロー /users/me/following?id=?
		v1.POST("users/me/following", FollowUser)
		v1.DELETE("users/me/following", UnfollowUser)
		// 自分がフォローしてるユーザー全部
		v1.GET("users/me/following", GetFollowingUser)
		// 自分を、フォローしてるユーザー全部
		v1.GET("users/me/followed", GetFollowedUser)
		// useridのuserがフォローしてるuserのツイート一覧 /users/:userid/followees/tweets
		v1.GET("tweets/following", GetFollowingUserTweets)
	}
	router.Run(":8080")
}
