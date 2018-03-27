package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"net/http"
	"os"
	"time"
)

func GenBody(c *gin.Context, world WorldState) {
	c.HTML(http.StatusOK, "index.tmpl.html", gin.H{"body": len(world.People), "people": world.People})
}
func genResponse(c *gin.Context, world WorldState) {
	genBody(c, world)
}

func genBody(c *gin.Context, world WorldState) {
	GenBody(c, world)
}

var theWorld WorldState

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	rand.Seed(time.Now())
	theWorld = WorldState{MeatLossFrac: 0.01, PerRoundLossFrac: 0.01, NewEntrantMeanAltruism: 10, NewEntrantMeanMeat: MAXMEAT / 10, People: make([]*Person, 0)}
	theWorld.initializeState()
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	//	router.GET("/", func(c *gin.Context) {
	//		c.String(http.StatusOK, "Hello World")
	//	})
	router.GET("/", func(c *gin.Context) {
		genResponse(c, theWorld)
		theWorld.updateState()
	})

	router.Run(":" + port)
}
