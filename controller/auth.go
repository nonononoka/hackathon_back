package controller

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
	"strings"
)

func authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/twitter/api/v1/hello" {
			ctx.Next()
			return
		}
		env := os.Getenv("ENV")
		log.Printf("authMiddleware")
		var opt option.ClientOption

		if env == "prod" {
			// 本番環境
			credentials, err := google.CredentialsFromJSON(ctx, []byte(os.Getenv("GOOGLE_CREDENTIALS_JSON")))
			if err != nil {
				log.Printf("error initializing app: %v\n", err)
				os.Exit(1)
			}
			opt = option.WithCredentials(credentials)
		} else {
			// local環境
			opt = option.WithCredentialsFile(os.Getenv("GOOGLE_CREDENTIALS_JSON"))
		}

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
