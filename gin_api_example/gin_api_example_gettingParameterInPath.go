package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func gettingWithParameterInPath(c *gin.Context) {
	name := c.Param("name")
	action := c.Param("action")
	log.Print("got http request")
	c.JSON(http.StatusOK, gin.H{"message": "hey",
		"name" : name,
		"action" : action,
		"status": http.StatusOK})
	
}
