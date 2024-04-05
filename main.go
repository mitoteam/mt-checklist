package main

import (
	"github.com/mitoteam/goappbase"
)

func main() {
	app := goappbase.NewAppBase()

	app.AppName = "MiTo Team Checklist"
	app.ExecutableName = "mt-checklist"
	app.LongDescription = `Checklists management system`

	app.Run()
}
