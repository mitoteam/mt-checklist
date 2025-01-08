package app

import (
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mt-checklist/model"
)

var App *goapp.AppBase

func InitApp() *goapp.AppBase {
	App = goapp.NewAppBase(defaultSettings)

	App.AppName = "MiTo Team Checklist"
	App.ExecutableName = "mt-checklist"
	App.LongDescription = `Checklists management system`

	App.PreRunF = DoPreRun
	App.PostRunF = DoPostRun

	return App
}

func DoPreRun() (err error) {
	// open database and migrate schema
	if err = goapp.DbSchema.Open(App.AppSettings.(*AppSettingsType).LogSql); err != nil {
		return err
	}

	//check if root user exists
	if err = model.InitializeRootUser(App.AppSettings.(*AppSettingsType).InitialRootPassword); err != nil {
		return err
	}

	return nil //no errors
}

func DoPostRun() error {
	goapp.DbSchema.Close()

	return nil //no errors
}
