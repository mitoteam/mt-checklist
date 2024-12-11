package web

import (
	"html/template"
	"net/http"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mitoteam/mt-checklist/app"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mttools"
	"github.com/mitoteam/mtweb"
)

func BuildWebRouter(r *gin.Engine) {
	r.StaticFileFS("/favicon.ico", "/favicon.ico", webAssetsHttpFS)

	//serve assets
	r.StaticFS("/assets", webAssetsHttpFS)

	//serve HTML from templates
	inc_templates := []string{"inc/base.html", "inc/header.html", "inc/footer.html"}

	render := multitemplate.NewRenderer()
	render.Add("admin_checklists", template.Must(template.ParseFS(templatesFS, append(inc_templates, "admin_checklists.html")...)))

	r.HTMLRender = render

	// no auth required routes
	r.GET("/logout", webLogout)

	r.GET("/sign-in", webDhtmlTemplate(PageLogin))
	r.POST("/sign-in", webDhtmlTemplate(PageLogin))

	// auth required routes
	authenticated_routes := r.Group("")
	authenticated_routes.Use(authMiddleware())
	authenticated_routes.
		GET("/", webDhtmlTemplate(PageDashboard))

	// Subgroup: admin role required routes
	admin_routes := authenticated_routes.Group("/admin")
	admin_routes.Use(adminRoleMiddleware())
	admin_routes.
		GET("/checklists", func(c *gin.Context) { c.HTML(http.StatusOK, "admin_checklists", buildRequestData(c)) })

	//EXPERIMENTS
	r.GET("/experiment", func(c *gin.Context) {
		c.Header("Content-Type", "text/html;charset=utf-8")
		c.String(http.StatusOK, mtweb.BuildExperimentHtml())
	})

	r.GET("/form", webDhtmlTemplate(PageFormExperiment))
	r.POST("/form", webDhtmlTemplate(PageFormExperiment))
}

// checks if user authenticated, redirects to /login if not (except for excludedPaths).
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		user := app.GetUser(mttools.AnyToInt64OrZero(session.Get("userID")))

		if user == nil {
			c.Redirect(http.StatusSeeOther, "/sign-in?destination="+c.Request.RequestURI)
			c.Abort() //stop other handlers
			return
		} else {
			c.Set("User", user)
		}

		// Call the next handler
		c.Next()
	}
}

func adminRoleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user *model.MtUser
		if v, ok := c.Get("User"); ok {
			user = v.(*model.MtUser)
		}

		if !user.HasRole(model.ROLE_ADMIN) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Call the next handler
		c.Next()
	}
}

// Prepares default set of gin.H data from context
func buildRequestData(c *gin.Context) gin.H {
	var user *model.MtUser
	if v, ok := c.Get("User"); ok {
		user = v.(*model.MtUser)
	}

	data := gin.H{
		"App": app.App,
	}

	if user != nil {
		data["UserID"] = user.ID
		data["User"] = user
	}

	return data
}

func webLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("userID")
	session.Save()

	c.Redirect(http.StatusFound, "/")
}

func webDhtmlTemplate(renderF func(*PageBuilder) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := NewPageBuilder(c)
		if renderF(p) {
			c.Header("Content-Type", "text/html;charset=utf-8")
			c.String(http.StatusOK, p.String())
		} else {
			c.Abort()
		}
	}
}
