package web

import (
	"net/http"

	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mttools"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := Session(r)

		user := model.LoadUser(mttools.AnyToInt64OrZero(session.Values["userID"]))
		//log.Printf("DBG AuthMiddleware User: %+v\n", user)

		if user == nil {
			url := mbr.Url(RootCtl.Login) + "?destination=" + r.RequestURI
			http.Redirect(w, r, url, http.StatusSeeOther)
			// do not call other handlers
		} else {
			if ctx := mbr.Context(r); ctx != nil {
				//fmt.Printf("DBG: mbr.Route %s user set\n", ctx.Route().Name())
				ctx.Set("User", user)
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		}
	})
}
