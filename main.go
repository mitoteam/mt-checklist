package main

import (
	_ "embed"

	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-checklist/app"
	"github.com/mitoteam/mt-checklist/web"
)

//go:embed LICENSE.md
var licenseString string

func main() {
	app := app.InitApp()
	app.SetHandler(mbr.Handler(web.RootCtl))

	app.License = licenseString

	app.Run()
}
