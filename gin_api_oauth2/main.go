package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions/cookie"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Credentials struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURIs []string `json:"redirect_uris"`
	AuthURI      string   `json:"auth_uri"`
	TokenURI     string   `json:"token_uri"`
}
type j struct {
	Web       *Credentials `json:"web"`
	Installed *Credentials `json:"installed"`
}

// User is a retrieved and authentiacted user.
type User struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Gender        string `json:"gender"`
}

var cred j
var conf *oauth2.Config
var state string
var store = cookie.NewStore([]byte("secret"))

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func init() {
	file, err := ioutil.ReadFile("/Users/michael/Downloads/test-clientid.google.json")
	if err != nil {
		log.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &cred)

	conf = &oauth2.Config{
		ClientID:     cred.Installed.ClientID,
		ClientSecret: cred.Installed.ClientSecret,
		RedirectURL:  "http://127.0.0.1:9090/google/auth",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}

	log.Println("Michael begin")
	log.Println(conf.ClientID)
	log.Println(conf.ClientSecret)
	log.Println(conf.Endpoint)
	log.Println("Michael end")
}

func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":       "hello gin " + strings.ToLower(c.Request.Method) + " method",
		"googleLogin": "/google/login",
	})
}

func getLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

func authHandler(c *gin.Context) {
	// Handle the exchange code to initiate a transport.
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	if retrievedState != c.Query("state") {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s with query %s", retrievedState, c.Query("state")))
		return
	}

	tok, err := conf.Exchange(oauth2.NoContext, c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	client := conf.Client(oauth2.NoContext, tok)
	email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer email.Body.Close()
	data, _ := ioutil.ReadAll(email.Body)
	log.Println("Email body: ", string(data))
	c.HTML(http.StatusOK, "googleSuccess.tmpl", gin.H{
		"state": retrievedState,
	})
}

func loginHandler(c *gin.Context) {
	state = randToken()
	session := sessions.Default(c)
	session.Set("state", state)
	session.Save()
	c.HTML(http.StatusOK, "googleLogin.tmpl", gin.H{
		"loginURL": getLoginURL(state), "state": state,
	})
	//c.Writer.Write([]byte("<html><title>Golang Google</title> <body> <a href='" + getLoginURL(state) +
	//	"'><button>Login with Google!</button> </a> </body></html>"))
}

func main() {
	router := gin.Default()
	router.Use(sessions.Sessions("mysession", store))
	router.Static("/statics", "./statics")
	router.LoadHTMLGlob("templates/*")
	router.StaticFile("/favicon.ico", "./favicon.ico")

	router.GET("/", indexHandler)
	googleLoginRouter := router.Group("/google")
	{
		googleLoginRouter.GET("/login", loginHandler)
		googleLoginRouter.GET("/auth", authHandler)
	}

	router.Run("127.0.0.1:9090")
}
