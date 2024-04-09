package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitoteam/goappbase"
	"github.com/mitoteam/mt-checklist/app"
)

func main() {
	application := goappbase.NewAppBase()

	application.AppName = "MiTo Team Checklist"
	application.ExecutableName = "mt-checklist"
	application.LongDescription = `Checklists management system`

	application.AppSettings = &app.AppSettingsType{}

	application.BuildWebRouterF = func(r *gin.Engine) {
		r.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "index.html") })
	}

	application.Run()
}
