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

	application.PreRunF = doPreRun
	application.PostRunF = doPostRun

	application.Run()
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
