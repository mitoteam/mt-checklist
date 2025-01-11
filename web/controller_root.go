package web

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mtweb"
)

type RootController struct {
	mbr.ControllerBase
}

var RootCtl *RootController

func init() {
	RootCtl = &RootController{}

	//using chi middlewares
	RootCtl.With(middleware.Recoverer)
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
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			cards_list := dhtmlbs.NewCardList().Class("row-cols-lg-2").Add(
				dhtmlbs.NewCard().Header(mtweb.Icon(iconChecklist).Label("Active checklists")).
					Body(c.renderActiveChecklists()),
			).Add(
				dhtmlbs.NewCard().Header(mtweb.Icon("user-check").Label("My issues")).Body("Some content"),
			).Add(
				dhtmlbs.NewCard().Header(mtweb.Icon("chart-simple").Label("Statistics")).Body(c.renderStatistics()),
			)

			if p.User().IsAdmin() {
				cards_list.Add(
					dhtmlbs.NewCard().Header(mtweb.Icon("cog").Label("System management")).Body(c.renderManagement()),
				).Add(
					dhtmlbs.NewCard().Header(mtweb.Icon("vial").Label("Experiment")).
						Body(
							dhtml.Div().Text("Html renderer ").
								Append(dhtml.NewLink(mbr.Url(RootCtl.Experiment)).Label("experiment")).Text(" link."),
						).
						Body(
							dhtml.Div().Text("Confirm link ").Append(dhtml.NewConfirmLink("/experiment", "Are you sure?").Label("experiment")),
						).
						Body(dhtml.Div().Text("Forms ").Append(dhtml.NewLink("/form").Label("experiment")).Text(" link.")),
				)
			}

			p.Main(cards_list)
			return nil
		}),
	}

	route.With(AuthMiddleware) //for home page only
	return route
}

func (c *RootController) renderStatistics() (out dhtml.HtmlPiece) {
	out.Append(
		dhtml.RenderValue(mtweb.Icon(iconUser).Label("Users"), goapp.CountOL[model.User]()),
		dhtml.RenderValueE(mtweb.Icon(iconChecklist).Label("Checklists"), goapp.CountOL[model.Checklist](), "no checklists created"),
	)

	out.Append(
		dhtml.RenderValueE(mtweb.Icon(iconTemplate).Label("Templates"), goapp.CountOL[model.Template](), "no templates created"),
	)

	return out
}

func (c *RootController) renderManagement() (out dhtml.HtmlPiece) {
	out.Append(dhtml.Div().Append(
		dhtml.NewLink(mbr.Url(AdminCtl.Users)).Label(mtweb.Icon(iconUser).Label("Users")),
	)).Append(
		dhtml.Div().Append(
			dhtml.NewLink(mbr.Url(AdminCtl.Templates)).Label(mtweb.Icon(iconTemplate).Label("Templates")),
		),
	).Append(
		dhtml.Div().Append(
			dhtml.Div().Append(
				dhtml.NewLink(mbr.Url(AdminCtl.Checklists)).Label(mtweb.Icon(iconChecklist).Label("Checklists")),
			).Append(" (administer)"),
		),
	)

	return out
}

func (c *RootController) renderActiveChecklists() (out dhtml.HtmlPiece) {
	list := goapp.LoadOL[model.Checklist]()

	div := dhtml.Div().Class("d-flex flex-wrap gap-2")

	for _, cl := range list {
		if cl.IsActive() {
			div.Append(
				mtweb.NewIconBtn(mbr.Url(ChecklistCtl.Checklist, "checklist_id", cl.ID), iconChecklist, cl.Name),
			)
		}
	}

	if !div.HasChildren() {
		div.Append(dhtml.EmptyLabel("no active checklists"))
	}

	out.Append(div)

	return out
}

func (c *RootController) Login() mbr.Route {
	return mbr.Route{
		PathPattern: "/login",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			user := p.User()

			if user == nil {
				p.Main(formLogin.Render(p.FormContext()))
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
			delete(session.Values, sessionIdField)
			err := session.Save(ctx.Request(), ctx.Writer())

			if err != nil {
				return err
			}

			ctx.RedirectRoute(http.StatusFound, RootCtl.Home)
			return nil
		},
	}
}

func (c *RootController) MyAccount() mbr.Route {
	route := mbr.Route{
		PathPattern: "/account",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			p.Main(formMyAccount.Render(p.FormContext()))
			return nil
		}),
	}

	route.With(AuthMiddleware)

	return route
}

func (c *RootController) Experiment() mbr.Route {
	return mbr.Route{
		PathPattern: "/experiment",
		HandleF:     func(ctx *mbr.MbrContext) any { return mtweb.BuildExperimentHtml() },
	}
}

func (c *RootController) TestForm() mbr.Route {
	return mbr.Route{
		PathPattern: "/form",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			p.Title("Form!").Main(mtweb.ExperimentFormHandler.Render(p.FormContext()))
			return nil
		}),
	}
}
