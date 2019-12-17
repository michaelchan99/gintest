package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func getting(c *gin.Context) {
	firstname := c.DefaultQuery("firstname", "")
	lastname := c.Query("lastname")
	log.Print("got http request")
	c.JSON(http.StatusOK, gin.H{"message": "hey",
		"firstname" : firstname,
		"lastname" : lastname,
		"status": http.StatusOK})

}