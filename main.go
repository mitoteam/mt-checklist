package main

import (
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-checklist/app"
	"github.com/mitoteam/mt-checklist/web"
)

func main() {
	app := app.InitApp()
	app.SetHandler(mbr.Handler(web.RootCtl))

	app.Run()
}
