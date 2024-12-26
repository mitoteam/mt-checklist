package web

import (
	"net/http"

	"github.com/mitoteam/mbr"
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

	route.With(AuthMiddleware) //for home page only
	return route
}

func (c *RootController) Login() mbr.Route {
	return mbr.Route{
		PathPattern: "/login",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			user := p.User()

			if user == nil {
				p.Form(Forms.Login)
			} else {
				p.Main("Already authenticated")
			}

			return nil
		}),
	}
}

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
