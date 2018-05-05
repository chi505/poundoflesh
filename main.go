package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"net/http"
	"os"
	"states"
)

func GenBody(c *gin.Context, world WorldState) {
	c.HTML(http.StatusOK, "index.tmpl.html", gin.H{"body": world.Count, "people": world.People})
}
func genResponse(c *gin.Context, world WorldState) {
	GenBody(c, world)
}

var theWorld WorldState

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	theWorld = WorldState{Count: 0, Params: PoundOFleshParams{MeatLossFrac: 0.2, PerRoundLossFrac: 0.05, NewEntrantMeanAltruism: 10, NewEntrantMeanMeat: int(MAXMEAT / 10)}, Assets: TextAssets{Organs: map[string][]MeatData{}}, People: make([]*Person, 0), PersonSpec: map[string]MeatSpec{}}
	theWorld.initializeState()
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		genResponse(c, theWorld)
		theWorld.updateState()
	})

	router.Run(":" + port)
}
