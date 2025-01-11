package web

import (
	"net/http"

	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mttools"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//log.Println("DBG AuthMiddleware " + r.RequestURI)

		ctx := mbr.Context(r)
		mttools.AssertNotNil(ctx, "not in MbrContext")

		session := Session(r)
		sessionId := ""
		if sessionIdValue, ok := session.Values[sessionIdField]; ok {
			sessionId = mttools.AnyToString(sessionIdValue)
		}

		goapp.PreQuery[model.User]().Where("session_id = ?", sessionId)
		user := goapp.FirstO[model.User]()

		if user == nil {
			url := mbr.Url(RootCtl.Login, "destination", r.RequestURI)
			http.Redirect(w, r, url, http.StatusSeeOther)
			// do not call other handlers
		} else {
			//fmt.Printf("DBG: mbr.Route %s user set\n", ctx.Route().Name())
			ctx.Set("User", user)

			// Call the next handler
			next.ServeHTTP(w, r)
		}
	})
}

func AdminRoleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//log.Println("DBG AdminRoleMiddleware " + r.RequestURI)

		ctx := mbr.Context(r)
		mttools.AssertNotNil(ctx, "not in MbrContext")

		var user *model.User
		if userValue, ok := ctx.GetOk("User"); ok {
			user = userValue.(*model.User)
		}

		if user != nil && user.HasRole(model.USER_ROLE_ADMIN) {
			// Call the next handler
			next.ServeHTTP(w, r)
		} else {
			ctx.ErrorWithCode(http.StatusUnauthorized, "Not authorized")
		}
	})
}
