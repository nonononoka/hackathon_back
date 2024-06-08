package controller

import (
	"back/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTags(ctx *gin.Context) {
	tags, error := usecase.GetTags()

	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, tags)
}
