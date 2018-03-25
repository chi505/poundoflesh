package main

import (
	"dummy"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"net/http"
	"os"
	"strconv"
)

var count = &dummy.MutableInt{0}

var People []Person

type MutableInt struct {
	Value int
}

func GenBody(c *gin.Context) {

	count.Value++
	c.HTML(http.StatusOK, "index.tmpl.html", gin.H{"body": strconv.Itoa(count.Value), "people": People})
}
func genResponse(c *gin.Context) {
	updateState()
	genBody(c)
}

func genBody(c *gin.Context) {
	dummy.GenBody(c, count, People)
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
