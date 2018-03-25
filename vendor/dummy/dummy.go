package dummy

import (
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"net/http"
	"strconv"
)

type MutableInt struct {
	Value int
}

func GenBody(c *gin.Context, count *MutableInt) {

	count.Value++
	c.HTML(http.StatusOK, "index.tmpl.html", gin.H{"body": strconv.Itoa(count.Value)})
}
