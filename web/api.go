package web

import (
	"github.com/mitoteam/goappbase"
	"github.com/mitoteam/mt-checklist/app"
)

// builds API routing and handlers for goappbase
func BuildWebApiRouter(application *goappbase.AppBase) {
	application.WebApiPathPrefix = "/api"
	application.WebApiEnableGet = !app.Settings.Production // in DEV mode only

	application.
		ApiHandler("/admin/checklists", api_AdminChecklists).
		ApiHandler("/ping", api_HealthCheck)
}

// returns true if user is logged in
func apiCheckSession(r *goappbase.ApiRequest) bool {
	if r.SessionGet("userID") == nil {
		r.SetErrorStatus("Auth Required")
		return false
	} else {
		return true
	}
}

func api_HealthCheck(r *goappbase.ApiRequest) error {
	r.SetOutData("session", apiCheckSession(r))
	r.SetOkStatus("API works: " + app.App.AppName)

	return nil
}
