package controller

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
	"strings"
)

func authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Printf("authMiddleware")
		opt := option.WithCredentialsFile(os.Getenv("GOOGLE_CREDENTIALS_JSON"))

		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			log.Printf("error initializing app: %v\n", err)
			os.Exit(1)
		}
		auth, err := app.Auth(ctx)
		if err != nil {
			log.Printf("error initializing auth: %v\n", err)
			os.Exit(1)
		}

		authHandler := ctx.Request.Header.Get("Authorization")
		idToken := strings.Replace(authHandler, "Bearer ", "", 1)
		log.Printf(idToken)
		token, err := auth.VerifyIDToken(ctx, idToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"Error": err.Error(),
			})
			return
		}
		if token.UID == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"msg":    "Unauthorized",
				"detail": "invalid token",
			})
			ctx.Abort()
			return
		}
		ctx.Set("token", token)
		ctx.Next()
	}
}
