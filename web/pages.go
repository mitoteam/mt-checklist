package web

import (
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mtweb"
)

func PageFormExperiment(p *PageBuilder) bool {
	formOut := dhtml.FormManager.RenderForm("test_form", p.FormContext())

	if !formOut.IsEmpty() {
		p.Title("Form!").Main(formOut)
		return true
	}

	return false
}

func PageDashboard(p *PageBuilder) bool {
	cards_list := mtweb.NewCardList().Add(
		mtweb.NewCard().Header(mtweb.Icon("vial").Label("Experiment")).
			Body(dhtml.Div().Text("Html renderer ").Append(dhtml.NewLink("/experiment").Label("experiment")).Text(" link.")).
			Body(
				dhtml.Div().Text("Confirm link ").Append(dhtml.NewConfirmLink("/experiment", "Are you sure?").Label("experiment")),
			).
			Body(dhtml.Div().Text("Forms ").Append(dhtml.NewLink("/form").Label("experiment")).Text(" link.")),
	).Add(
		mtweb.NewCard().Header(mtweb.Icon("list-check").Label("Active checklists")).
			Body("Some content"),
	).Add(
		mtweb.NewCard().Header(mtweb.Icon("user-check").Label("My issues")).
			Body("Some content"),
	).Add(
		mtweb.NewCard().Header(mtweb.Icon("chart-simple").Label("Statistics")).Body(renderStatistics()),
	).Add(
		mtweb.NewCard().Header(mtweb.Icon("cog").Label("System management")).
			Body(dhtml.Div().Append(
				dhtml.NewLink("/admin/users").Label(mtweb.Icon("user").Label("Users")),
			)).
			Body(dhtml.Div().Append(
				dhtml.Div().Append(
					dhtml.NewLink("/admin/checklists").Label(mtweb.Icon("list-check").Label("Checklists")),
				).Append(" (administration)"),
			)),
	)

	p.Main(cards_list)
	return true
}

func renderStatistics() (out dhtml.HtmlPiece) {
	out.Append(dhtml.RenderValue("User count", len(model.GetUsersList())))
	out.Append(dhtml.RenderValueE("Checklist count", len(model.GetChecklistsList()), "no checklists created"))

	return out
}

func PageLogin(p *PageBuilder) bool {
	p.Title("Sign In")

	session := p.GetSession()

	if userID, ok := session.Get("userID").(int64); ok && userID > 0 {
		p.Main("Already authenticated")
	} else {
		fc := p.FormContext().
			SetParam("Session", session).
			SetRedirect(p.GetGinContext().DefaultQuery("destination", "/"))

		formOut := dhtml.FormManager.RenderForm("login", fc)

		if formOut.IsEmpty() {
			return false
		} else {
			p.Main(formOut)
		}
	}

	return true
}
