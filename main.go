package main

import (
	"github.com/mitoteam/goappbase"
	"github.com/mitoteam/mt-checklist/app"
	"github.com/mitoteam/mt-checklist/web"
)

func main() {
	application := goappbase.NewAppBase(app.Settings)

	application.AppName = "MiTo Team Checklist"
	application.ExecutableName = "mt-checklist"
	application.LongDescription = `Checklists management system`

	application.BuildWebRouterF = web.BuildWebRouter

	application.Run()
}
