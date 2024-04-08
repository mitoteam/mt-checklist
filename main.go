package main

import (
	"github.com/mitoteam/goappbase"
	"github.com/mitoteam/mt-checklist/app"
)

func main() {
	application := goappbase.NewAppBase()

	application.AppName = "MiTo Team Checklist"
	application.ExecutableName = "mt-checklist"
	application.LongDescription = `Checklists management system`

	application.AppSettings = &app.AppSettingsType{}

	application.Run()
}
