package web

import (
	"html/template"
	"net/http"
	"slices"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mitoteam/mt-checklist/app"
)

func BuildWebRouter(r *gin.Engine) {
	r.StaticFileFS("/favicon.ico", "/favicon.ico", webAssetsHttpFS)

	//serve assets
	r.StaticFS("/assets", webAssetsHttpFS)

	//serve HTML from templates
	inc_templates := []string{"inc/base.html", "inc/header.html", "inc/footer.html"}

	render := multitemplate.NewRenderer()
	render.Add("index", template.Must(template.ParseFS(templatesFS, append(inc_templates, "index.html")...)))
	render.Add("login_form", template.Must(template.ParseFS(templatesFS, append(inc_templates, "login.html")...)))
	render.Add("checklist", template.Must(template.ParseFS(templatesFS, append(inc_templates, "checklist.html")...)))

	r.HTMLRender = render

	r.Use(authMiddleware([]string{"/login", "/logout"}))
	r.GET("/", webIndex)
	r.GET("/logout", webLogout)
	r.POST("/login", webLoginPost) //login form handler
}

// checks if user authenticated, redirects to /login if not (except for excludedPaths).
func authMiddleware(excludedPaths []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if slices.Contains(excludedPaths, c.FullPath()) {
			//no auth required for route
			c.Next()
			return
		}

		session := sessions.Default(c)

		if session.Get("userID") == nil {
			webLoginForm(c)
			c.Abort() //stop other handlers
			return
		}

		// Call the next handler
		c.Next()
	}
}

func buildRequestData(c *gin.Context) gin.H {
	session := sessions.Default(c)

	data := gin.H{
		"App":    app.App,
		"UserID": session.Get("userID"),
	}

	return data
}

func webIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index", buildRequestData(c))
}

func webLoginForm(c *gin.Context) {
	c.HTML(http.StatusOK, "login_form", buildRequestData(c))
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
