package main

import (
	"dummy"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"os"
)

var count = &dummy.MutableInt{0}

func genResponse(c *gin.Context) {
	updateState()
	genBody(c)
}

func genBody(c *gin.Context) {
	dummy.GenBody(c, count)
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	initializeState()
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	//	router.GET("/", func(c *gin.Context) {
	//		c.String(http.StatusOK, "Hello World")
	//	})
	router.GET("/", func(c *gin.Context) {
		genResponse(c)
	})

	router.Run(":" + port)
}
