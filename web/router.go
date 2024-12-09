package web

import (
	"html/template"
	"net/http"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/mt-checklist/app"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mtweb"
)

func BuildWebRouter(r *gin.Engine) {
	r.StaticFileFS("/favicon.ico", "/favicon.ico", webAssetsHttpFS)

	//serve assets
	r.StaticFS("/assets", webAssetsHttpFS)

	//serve HTML from templates
	inc_templates := []string{"inc/base.html", "inc/header.html", "inc/footer.html"}

	render := multitemplate.NewRenderer()
	render.Add("placeholder", template.Must(template.ParseFS(templatesFS, append(inc_templates, "placeholder.html")...)))
	render.Add("dashboard", template.Must(template.ParseFS(templatesFS, append(inc_templates, "dashboard.html")...)))
	render.Add("login_form", template.Must(template.ParseFS(templatesFS, append(inc_templates, "login.html")...)))
	render.Add("admin_checklists", template.Must(template.ParseFS(templatesFS, append(inc_templates, "admin_checklists.html")...)))
	render.Add("checklist", template.Must(template.ParseFS(templatesFS, append(inc_templates, "checklist.html")...)))

	r.HTMLRender = render

	// no auth required routes
	r.GET("/experiment", webExperiment)
	r.GET("/logout", webLogout)
	r.POST("/login", webLoginPost) //login form handler

	// auth required routes
	g_auth := r.Group("")

	g_auth.Use(authMiddleware()).
		GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "dashboard", buildRequestData(c)) })

	// admin role required routes
	g_auth.Group("/admin").
		Use(adminRoleMiddleware()).
		GET("/checklists", func(c *gin.Context) { c.HTML(http.StatusOK, "admin_checklists", buildRequestData(c)) })

	r.GET("/form", webPlaceholder("Form!", func(c *gin.Context) *dhtml.HtmlPiece {
		return dhtml.RenderForm(mtweb.GetTestForm(), c.Writer, c.Request)
	}))
	r.POST("/form", webPlaceholder("Form!", func(c *gin.Context) *dhtml.HtmlPiece {
		return dhtml.RenderForm(mtweb.GetTestForm(), c.Writer, c.Request)
	}))
}

// checks if user authenticated, redirects to /login if not (except for excludedPaths).
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get("userID") == nil {
			c.HTML(http.StatusOK, "login_form", buildRequestData(c))
			c.Abort() //stop other handlers
			return
		}

		// Call the next handler
		c.Next()
	}
}

func adminRoleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		user := app.GetUser(session.Get("userID").(uint))

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
	session := sessions.Default(c)

	user_id := session.Get("userID")
	data := gin.H{
		"App":    app.App,
		"UserID": user_id,
	}

	if user_id != nil {
		data["User"] = app.GetUser(user_id.(uint))
	}

	return data
}

// Login form POST handler
func webLoginPost(c *gin.Context) {
	session := sessions.Default(c)

	var errMessage string

	username := c.PostForm("username")

	if username == "" {
		errMessage += "Username not given\n"
	}

	password := c.PostForm("password")

	if password == "" {
		errMessage += "Password not given\n"
	}

	if errMessage == "" {
		user := app.AuthorizeUser(username, password)

		if user != nil {
			session.Set("userID", user.ID)
			session.Save()

			c.Redirect(http.StatusFound, "/")
		} else {
			session.Delete("userID")
			session.Save()

			errMessage = "User not found or wrong password given"
		}
	}

	if errMessage != "" {
		errMessage = "<pre>" + errMessage + "</pre><div><a href=\"/\">Main Page</a></div>"

		c.Header("Content-Type", "text/html;charset=utf-8")
		c.String(http.StatusUnauthorized, errMessage)
		c.Abort()
		return
	}
}

func webLogout(c *gin.Context) {
	session := sessions.Default(c)

	session.Delete("userID")
	session.Save()

	c.Redirect(http.StatusFound, "/")
}

// Builds handler function for placeholder.html template
func webPlaceholder(page_title string, builderF func(*gin.Context) *dhtml.HtmlPiece) gin.HandlerFunc {
	return func(c *gin.Context) {
		responseHtml := builderF(c)

		if !responseHtml.IsEmpty() {
			data := buildRequestData(c)
			data["Title"] = page_title
			data["Content"] = template.HTML(
				dhtml.Piece(
					mtweb.NewCard().
						Header(dhtml.Span().Class("fs-5").Text(page_title)).
						Body(builderF(c)),
				).String(),
			)

			c.HTML(http.StatusOK, "placeholder", data)
		}
	}
}

// Handler for /experiment path
func webExperiment(c *gin.Context) {
	c.Header("Content-Type", "text/html;charset=utf-8")

	c.String(http.StatusOK, mtweb.BuildExperimentHtml())
}
