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
		v1.GET("tweets", GetTweets)
		v1.POST("tweet", PostTweet)
		v1.GET("tags", GetTags)
		// 特定のuserのツイートを取得 /users/:id/tweets
		v1.GET("users/:id/tweets", GetUserTweets)
		// 全部のuserを取得 /users
		v1.GET("users", GetUsers)
		// :idのuserをフォロー /followees/:userid post
		v1.POST("friendships/:id/follow", FollowUser)
		// :idがフォローしてるuserを取得 /users/:userid/followees get (:useridにフォローされてる人たち)
		v1.GET("friendships/follow", GetFollowingUser)
		// useridのuserがフォローしてるuserのツイート一覧 /users/:userid/followees/tweets
		v1.GET("friendships/tweets", GetFollowingUserTweets)
		// 特定ツイートのいいね /:idのツイートをファボする。
		v1.POST("tweets/:id/favorites", PostLike)
		//// 特定userがいいねしてるツイートを取得する /users/:userid/tweets/favorites
		//v1.GET("users/:id/likes-tweets", GetLikeTweets)
		// 特定ツイートにリプライする。
		v1.POST("tweets/:id/reply", PostReply)
	}
	router.Run(":8080")
}
