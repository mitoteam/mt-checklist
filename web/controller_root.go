package web

import (
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
	return mbr.Route{PathPattern: "/assets", StaticFS: webAssetsFS}
}

func (c *RootController) FavIcon() mbr.Route {
	return mbr.Route{PathPattern: "/favicon.ico", FileFromFS: "favicon.ico", StaticFS: webAssetsFS}
}

func (c *RootController) Home() mbr.Route {
	return mbr.Route{
		PathPattern: "/",
		HandleF:     PageBuilderRouteHandler(PageDashboard),
	}
}
