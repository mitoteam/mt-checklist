package web

import (
	"html/template"
	"net/http"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/mitoteam/mt-checklist/app"
)

func BuildWebRouter(r *gin.Engine) {
	r.StaticFileFS("/favicon.ico", "/favicon.ico", webAssetsHttpFS)

	//serve assets
	r.StaticFS("/assets", webAssetsHttpFS)

	//serve HTML from templates
	inc_templates := []string{"inc/base.html", "inc/footer.html"}

	render := multitemplate.NewRenderer()
	render.Add("index", template.Must(template.ParseFS(templatesFS, append(inc_templates, "index.html")...)))
	render.Add("login", template.Must(template.ParseFS(templatesFS, append(inc_templates, "login.html")...)))
	render.Add("checklist", template.Must(template.ParseFS(templatesFS, append(inc_templates, "checklist.html")...)))

	r.HTMLRender = render

	r.GET("/", webIndex)
	r.GET("/login", webLogin)
}

func webIndex(c *gin.Context) {
	//session := sessions.Default(c)

	data := gin.H{
		"App": app.App,
	}

	c.HTML(http.StatusOK, "index", data)
}

func webLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{})
}
