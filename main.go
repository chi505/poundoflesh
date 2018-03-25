package main

import (
	"dummy"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	count := dummy.MutableInt(0)
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	//	router.GET("/", func(c *gin.Context) {
	//		c.String(http.StatusOK, "Hello World")
	//	})
	router.GET("/", func(c *gin.Context) {
		dummy.GenBody(c, count)
	})

	router.Run(":" + port)
}
