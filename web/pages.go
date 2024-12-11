package web

import (
	"github.com/gin-contrib/sessions"
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/mtweb"
)

func PageDashboard(p *PageTemplate) bool {
	cards_list := mtweb.NewCardList().
		Add(
			mtweb.NewCard().Header(mtweb.Icon("vial").Label("Experiment")).
				Body(dhtml.Div().Text("Html renderer ").Append(dhtml.NewLink("/experiment").Label("experiment")).Text(" link.")).
				Body(dhtml.Div().Text("Forms ").Append(dhtml.NewLink("/form").Label("experiment")).Text(" link.")),
		).
		Add(
			mtweb.NewCard().Header(mtweb.Icon("list-check").Label("Active checklists")).
				Body("Some content"),
		).
		Add(
			mtweb.NewCard().Header(mtweb.Icon("user-check").Label("My issues")).
				Body("Some content"),
		).
		Add(
			mtweb.NewCard().Header(mtweb.Icon("chart-simple").Label("Statistics")).
				Body("Some content"),
		).
		Add(
			mtweb.NewCard().Header(mtweb.Icon("cog").Label("System management")).
				Body(
					dhtml.Div().Append(
						dhtml.NewLink("/admin/checklists").Label(mtweb.Icon("list-check").Label("Checklists")),
					).Append(" (administration)"),
				),
		)

	p.Main(cards_list)
	return true
}

func PageLogin(p *PageTemplate) bool {
	p.Title("Sign In")

	session := sessions.Default(p.GetContext())

	if userID, ok := session.Get("userID").(int64); ok && userID > 0 {
		p.Main("Already authenticated")
	} else {
		fc := FormContextFromGin(p.GetContext()).
			SetParam("Session", session).
			SetRedirect(p.GetContext().DefaultQuery("destination", "/"))

		formOut := dhtml.FormManager.RenderForm("login", fc)

		if formOut.IsEmpty() {
			return false
		} else {
			p.Main(formOut)
		}
	}

	return true
}
