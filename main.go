package main

import (
	"flag"
	"fmt"
	"gintest/googleAuth"
	"github.com/gin-contrib/sessions/cookie"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":       "hello gin " + strings.ToLower(c.Request.Method) + " method",
		"googleLogin": "/googleLogin",
	})
}

var store = cookie.NewStore([]byte("secret"))
var redirectURL, credFile string

func init() {
	bin := path.Base(os.Args[0])
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s
================
`, bin)
		flag.PrintDefaults()
	}
	flag.StringVar(&redirectURL, "redirect", "http://127.0.0.1:9090/googleAuth/success", "URL to be redirected to after authorization.")
	flag.StringVar(&credFile, "cred-file", "/Users/michael/Downloads/test-clientid.google.json", "Credential JSON file")
}

func main() {
	flag.Parse()
	router := gin.Default()

	scopes := []string{
		"https://www.googleapis.com/auth/userinfo.email",
		// You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
	}

	googleAuth.Setup(redirectURL, credFile, scopes)

	router.Use(sessions.Sessions("mysession", store))
	router.Static("/statics", "./statics")
	router.LoadHTMLGlob("templates/*")
	router.StaticFile("/favicon.ico", "./favicon.ico")

	router.GET("/", indexHandler)
	router.GET("/googleLogin", googleAuth.LoginHandler)

	private := router.Group("/googleAuth")
	{
		private.Use(googleAuth.AuthHandler())
		private.GET("/api", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"message": "Hello from private for groups"})
		})
		private.GET("/success", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "googleSuccess.tmpl", gin.H{
				"state": "Success"})
		})

		//log.Println("Email body: ", string(data))
		//ctx.HTML(http.StatusOK, "googleSuccess.tmpl", gin.H{
		//	"state": retrievedState,
		//}))
	}

	router.Run("127.0.0.1:9090")
}
