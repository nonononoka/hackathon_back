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
		v1.GET("tweet", GetTweets)
		v1.POST("tweet", PostTweet)
		v1.GET("tags", GetTags)
	}
	router.Run(":8080")
}
