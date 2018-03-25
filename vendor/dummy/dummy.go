package dummy

import "strconv"
import "github.com/gin-gonic/gin"

type MutableInt struct {
	value int
}

func genBody(c *gin.Context, count *MutableInt) {

	count.value++
	c.HTML(http.StatusOK, "index.tmpl.html", gin.H{"body": strconv.Itoa(count.value)})
}
