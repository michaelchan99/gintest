package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	ginglog "github.com/szuecs/gin-glog"
	ginoauth2 "github.com/zalando/gin-oauth2"
	"time"
)

// https://github.com/zalando/gin-oauth2

func main() {
	flag.Parse()
	router := gin.New()
	router.Use(ginglog.Logger(3 * time.Second))
	router.Use(ginoauth2.RequestLogger([]string{"uid"}, "data"))
	router.Use(gin.Recovery())
}
