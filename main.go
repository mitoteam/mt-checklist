package main

import (
	"github.com/mitoteam/goappbase"
	"github.com/mitoteam/mt-checklist/app"
	"github.com/mitoteam/mt-checklist/web"
)

func main() {
	app.App = goappbase.NewAppBase(app.DefaultSettings)

	app.App.AppName = "MiTo Team Checklist"
	app.App.ExecutableName = "mt-checklist"
	app.App.LongDescription = `Checklists management system`

	app.App.BuildWebRouterF = web.BuildWebRouter

	app.App.PreRunF = app.DoPreRun
	app.App.PostRunF = app.DoPostRun

	app.App.Run()
}
