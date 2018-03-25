package dummy

import (
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"net/http"
	"os"
	"strconv"
)

type MutableInt struct {
	value int
}

func genBody(c *gin.Context, count *MutableInt) {

	count.value++
	c.HTML(http.StatusOK, "index.tmpl.html", gin.H{"body": strconv.Itoa(count.value)})
}
