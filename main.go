package main

import (
	"github.com/mitoteam/goappbase"
	"github.com/mitoteam/mt-checklist/app"
	"github.com/mitoteam/mt-checklist/web"
)

func main() {
	settings := &app.AppSettingsType{}
	settings.WebserverPort = 15119

	application := goappbase.NewAppBase(settings)

	application.AppName = "MiTo Team Checklist"
	application.ExecutableName = "mt-checklist"
	application.LongDescription = `Checklists management system`

	application.BuildWebRouterF = web.BuildWebRouter

	application.Run()
}
