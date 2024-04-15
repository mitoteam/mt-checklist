package web

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitoteam/mt-checklist/app"
)

func BuildWebRouter(r *gin.Engine) {
	r.StaticFileFS("/favicon.ico", "/favicon.ico", webAssetsHttpFS)

	//serve assets
	r.StaticFS("/assets", webAssetsHttpFS)

	//serve HTML from templates
	t := template.Must(template.New("index").ParseFS(templatesFS, "*.html"))
	r.SetHTMLTemplate(t)
	r.GET("/", webIndex)
}

func webIndex(c *gin.Context) {
	//session := sessions.Default(c)

	data := gin.H{
		"App": app.App,
	}

	c.HTML(http.StatusOK, "index.html", data)
}
