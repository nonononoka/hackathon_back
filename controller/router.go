package controller

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
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
		},
	}))

	v1 := router.Group("todo/api/v1")
	{
		v1.GET("todos", todosGET)
		v1.POST("todos", todoPOST)
	}
	router.Run(":8080")
}

func todosGET(c *gin.Context) {
	todos := []Todo{
		{ID: 1, Title: "Complete homework"},
		{ID: 2, Title: "Buy groceries"},
		{ID: 3, Title: "Go for a run"},
	}
	c.JSON(http.StatusOK, gin.H{"todos": todos})
}
