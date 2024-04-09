package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitoteam/goappbase"
	"github.com/mitoteam/mt-checklist/app"
)

func main() {
	settings := &app.AppSettingsType{}
	settings.WebserverPort = 15119

	application := goappbase.NewAppBase(settings)

	application.AppName = "MiTo Team Checklist"
	application.ExecutableName = "mt-checklist"
	application.LongDescription = `Checklists management system`

	application.BuildWebRouterF = func(r *gin.Engine) {
		r.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "index.html") })
	}

	application.Run()
}
