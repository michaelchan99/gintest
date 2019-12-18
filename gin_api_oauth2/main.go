package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/zalando/gin-oauth2/google"
	"net/http"
)

// https://github.com/zalando/gin-oauth2

func main() {

	redirectURL := "http://127.0.0.1:8081/auth/"
	credFile := "./example/google/test-clientid.google.json" // you have to build your own
	scopes := []string{
		"https://www.googleapis.com/auth/userinfo.email",
		// You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
	}
	secret := []byte("secret") //
	sessionName := "goquestsession"

	flag.Parse()

	router := gin.Default()
	// init settings for google auth
	google.Setup(redirectURL, credFile, scopes, secret)
	router.Use(google.Session(sessionName))
	router.Use(gin.Recovery())

	router.GET("/login", google.LoginHandler)

	// protected url group
	private := router.Group("/auth")
	private.Use(google.Auth())
	private.GET("/", UserInfoHandler)
	private.GET("/api", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Hello from private for groups"})
	})

	router.Run(":8081")
}

func UserInfoHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"Hello": "from private", "user": ctx.MustGet("user").(google.User)})
}
