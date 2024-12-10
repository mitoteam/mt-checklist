package web

import "github.com/gin-gonic/gin"

func PageDashboard(c *gin.Context, p *PageTemplate) {
	p.Main("dashboard")
}
