package controller

import (
	"back/usecase"
	"database/sql"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type UserConfig struct {
	Name  string `json:"name"`
	Bio   string `json:"bio"`
	Image string `json:"image"`
}

func GetMe(ctx *gin.Context) {
	token := ctx.MustGet("token").(*auth.Token)

	userInfo, error := usecase.GetMe(token)
	log.Println(userInfo)
	if error != nil {
		if error == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, userInfo)
}

func PostMe(ctx *gin.Context) {
	token := ctx.MustGet("token").(*auth.Token)

	userInfo, error := usecase.PostMe(token)
	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}
	log.Println(userInfo)
	ctx.JSON(http.StatusCreated, userInfo)
}

func PutMe(ctx *gin.Context) {
	log.Println(ctx)
	token := ctx.MustGet("token").(*auth.Token)
	var userConfig UserConfig
	if err := ctx.BindJSON(&userConfig); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}
	log.Println("userCOnfig", userConfig)

	userInfo, error := usecase.PutMe(token, userConfig.Name, userConfig.Bio, userConfig.Image)
	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}
	log.Println(userInfo)
	ctx.JSON(http.StatusCreated, userInfo)
}

func GetUsers(ctx *gin.Context) {
	token := ctx.MustGet("token").(*auth.Token)
	userInfo, error := usecase.GetUsers(token)

	if error != nil {
		if error == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, userInfo)
}
