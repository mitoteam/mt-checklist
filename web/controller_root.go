package web

import (
	"errors"

	"github.com/mitoteam/mbr"
)

type RootController struct {
	mbr.ControllerBase
}

var Root *RootController

func init() {
	Root = &RootController{}
}

func (c *RootController) Assets() mbr.Route {
	return mbr.Route{
		PathPattern: "/assets",
		StaticFS:    webAssetsFS,
	}
}

func (c *RootController) FavIcon() mbr.Route {
	return mbr.Route{
		PathPattern: "/favicon.ico",
		FileFromFS:  "favicon.ico",
		StaticFS:    webAssetsFS,
	}
}

func (c *RootController) Home() mbr.Route {
	return mbr.Route{
		PathPattern: "/",
		HandleF: func(ctx *mbr.MbrContext) any {
			p := NewPageBuilder(nil)
			if PageDashboard(p) {
				ctx.Request().Header.Add("Content-Type", "text/html;charset=utf-8")
				return p.String()
			} else {
				return errors.New("error")
			}
		},
	}
}
