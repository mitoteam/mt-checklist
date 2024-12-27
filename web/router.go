package web

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mttools"
)

func BuildWebRouter(r *gin.Engine) {

	// auth required routes
	authenticated_routes := r.Group("")
	authenticated_routes.Use(authMiddlewareOLD())

	// Subgroup: admin role required routes
	admin_routes := authenticated_routes.Group("/admin")
	admin_routes.Use(adminRoleMiddleware())
}

// checks if user authenticated, redirects to /login if not (except for excludedPaths).
func authMiddlewareOLD() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		user := model.LoadUser(mttools.AnyToInt64OrZero(session.Get("userID")))

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
		var user *model.User
		if v, ok := c.Get("User"); ok {
			user = v.(*model.User)
		}

		if !user.HasRole(model.USER_ROLE_ADMIN) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Call the next handler
		c.Next()
	}
}

func webPageBuilder(renderF func(*PageBuilderOLD) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := NewPageBuilderOLD(c)

		if renderF(p) {
			c.Header("Content-Type", "text/html;charset=utf-8")
			c.String(http.StatusOK, p.String())
		} else {
			c.Abort()
		}
	}
}
