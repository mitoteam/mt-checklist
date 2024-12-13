package app

import (
	"github.com/mitoteam/goappbase"
	"github.com/mitoteam/mt-checklist/model"
)

var App *goappbase.AppBase

func InitApp() *goappbase.AppBase {
	App = goappbase.NewAppBase(defaultSettings)

	App.AppName = "MiTo Team Checklist"
	App.ExecutableName = "mt-checklist"
	App.LongDescription = `Checklists management system`

	App.PreRunF = DoPreRun
	App.PostRunF = DoPostRun

	return App
}

func DoPreRun() (err error) {
	// open database and migrate schema
	if err = goappbase.DbSchema.Open(); err != nil {
		return err
	}

	//check if root user exists
	if err = model.InitializeRootUser(App.AppSettings.(*AppSettingsType).InitialRootPassword); err != nil {
		return err
	}

	return nil //no errors
}

func DoPostRun() error {
	goappbase.DbSchema.Close()

	return nil //no errors
}
