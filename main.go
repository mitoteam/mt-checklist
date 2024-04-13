package main

import (
	"github.com/mitoteam/goappbase"
	"github.com/mitoteam/mt-checklist/app"
	"github.com/mitoteam/mt-checklist/web"
)

func main() {
	app.App = goappbase.NewAppBase(app.Settings)

	app.App.AppName = "MiTo Team Checklist"
	app.App.ExecutableName = "mt-checklist"
	app.App.LongDescription = `Checklists management system`

	app.App.BuildWebRouterF = web.BuildWebRouter

	app.App.PreRunF = doPreRun
	app.App.PostRunF = doPostRun

	app.App.Run()
}

func doPreRun() error {
	var err error

	// open database and migrate schema
	if app.Db, err = goappbase.DbSchema.Open(); err != nil {
		return err
	}

	return nil //no errors
}

func doPostRun() error {
	goappbase.DbSchema.Close()

	return nil //no errors
}
