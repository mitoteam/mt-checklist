package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BuildWebRouter(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "index.html") })
}
