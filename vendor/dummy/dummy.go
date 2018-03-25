package dummy

import "strconv"

func genBody(c *gin.Context, count *int) {

	count = count + 1
	c.HTML(http.StatusOK, "index.tmpl.html", gin.H{"body": strconv.Itoa(count)})
}
