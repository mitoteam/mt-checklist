package web

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/mitoteam/mt-checklist/app"
)

var sessionStore *sessions.CookieStore
var sessionName = "mt-checklist"
var sessionIdField = "sessionId"

func Session(r *http.Request) *sessions.Session {
	if sessionStore == nil {
		sessionStore = sessions.NewCookieStore([]byte(app.App.AppSettings.(*app.AppSettingsType).WebserverCookieSecret))
	}

	session, _ := sessionStore.Get(r, sessionName)

	return session
}
