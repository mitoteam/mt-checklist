package main

import (
	"github.com/mitoteam/mt-checklist/app"
	"github.com/mitoteam/mt-checklist/web"
)

func main() {
	app := app.InitApp()
	app.BuildWebRouterF = web.BuildWebRouter

	app.Run()
}
