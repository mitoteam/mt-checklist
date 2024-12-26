package web

import (
	"net/http"

	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mttools"
	"github.com/mitoteam/mtweb"
)

type RootController struct {
	mbr.ControllerBase
}

var RootCtl *RootController

func init() {
	RootCtl = &RootController{}
}

func (c *RootController) Assets() mbr.Route {
	return mbr.Route{PathPattern: "/assets", StaticFS: webAssetsFS}
}

func (c *RootController) FavIcon() mbr.Route {
	return mbr.Route{PathPattern: "/favicon.ico", FileFromFS: "favicon.ico", StaticFS: webAssetsFS}
}

func (c *RootController) Home() mbr.Route {
	route := mbr.Route{
		PathPattern: "/",
		HandleF:     PageBuilderRouteHandler(PageDashboard),
	}

	route.With(AuthMiddleware)
	return route
}

func (c *RootController) Login() mbr.Route {
	return mbr.Route{
		PathPattern: "/login",
		HandleF: func(ctx *mbr.MbrContext) any {
			return "Login Form"
		},
	}
}

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

// func PageLogin(p *PageBuilder) bool {
// 	p.Title("Sign In")

// 	session := p.GetSession()

// 	if userID, ok := session.Get("userID").(int64); ok && userID > 0 {
// 		user := goapp.LoadO[model.User](userID)

// 		if user == nil {
// 			//Session user not found, restart session
// 			session.Clear()
// 			session.Save()

// 			p.GetGinContext().Redirect(http.StatusSeeOther, "/")
// 			return false
// 		} else {
// 			p.Main("Already authenticated")
// 		}
// 	} else {
// 		fc := p.FormContext().
// 			SetParam("Session", session).
// 			SetRedirect(p.GetGinContext().DefaultQuery("destination", "/"))

// 		formOut := Forms.Login.Render(fc)

// 		if formOut.IsEmpty() {
// 			return false
// 		} else {
// 			p.Main(formOut)
// 		}
// 	}

// 	return true
// }

func (c *RootController) Logout() mbr.Route {
	return mbr.Route{
		PathPattern: "/logout",
		HandleF: func(ctx *mbr.MbrContext) any {
			session := Session(ctx.Request())
			delete(session.Values, "userID")
			err := session.Save(ctx.Request(), ctx.Writer())

			if err != nil {
				return err
			}

			ctx.RedirectUrl(http.StatusFound, RootCtl.Home)
			return nil
		},
	}
}

func (c *RootController) Experiment() mbr.Route {
	return mbr.Route{
		PathPattern: "/experiment",
		HandleF:     func(ctx *mbr.MbrContext) any { return mtweb.BuildExperimentHtml() },
	}
}
