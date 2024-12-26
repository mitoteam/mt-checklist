package web

import (
	"net/http"

	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mtweb"
)

const (
	iconTemplate  = "pen-ruler"
	iconUser      = "person"
	iconChecklist = "list-check"
)

func PageFormExperiment(p *PageBuilderOLD) bool {
	formOut := mtweb.ExperimentFormHandler.Render(p.FormContext())

	if !formOut.IsEmpty() {
		p.Title("Form!").Main(formOut)
		return true
	}

	return false
}

func PageDashboard(p *PageBuilder) any {
	cards_list := mtweb.NewCardList().Add(
		mtweb.NewCard().Header(mtweb.Icon("vial").Label("Experiment")).
			Body(
				dhtml.Div().Text("Html renderer ").
					Append(dhtml.NewLink(mbr.Url(RootCtl.Experiment)).Label("experiment")).Text(" link."),
			).
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
	)

	if p.User().IsAdmin() {
		cards_list.Add(
			mtweb.NewCard().Header(mtweb.Icon("cog").Label("System management")).
				Body(dhtml.Div().Append(
					dhtml.NewLink("/admin/users").Label(mtweb.Icon(iconUser).Label("Users")),
				)).
				Body(dhtml.Div().Append(
					dhtml.NewLink("/admin/templates").Label(mtweb.Icon(iconTemplate).Label("Templates")),
				)).
				Body(dhtml.Div().Append(
					dhtml.Div().Append(
						dhtml.NewLink("/admin/checklists").Label(mtweb.Icon(iconChecklist).Label("Checklists")),
					).Append(" (administration)"),
				)),
		)
	}

	p.Main(cards_list)
	return nil
}

func renderStatistics() (out dhtml.HtmlPiece) {
	out.Append(
		dhtml.RenderValue(mtweb.Icon(iconUser).Label("Users"), goapp.CountOL[model.User]()),
		dhtml.RenderValueE(mtweb.Icon(iconChecklist).Label("Checklists"), goapp.CountOL[model.Checklist](), "no checklists created"),
	)

	goapp.PreQuery[model.Checklist]().Where("is_active = ?", true)
	out.Append(
		dhtml.RenderValueE(
			mtweb.Icon("flag").Label("Active checklists"), goapp.CountOL[model.Checklist](), "no active checklists",
		),
	)

	out.Append(
		dhtml.RenderValueE(mtweb.Icon(iconTemplate).Label("Templates"), goapp.CountOL[model.ChecklistTemplate](), "no templates created"),
	)

	return out
}

func PageLoginOLD(p *PageBuilderOLD) bool {
	p.Title("Sign In")

	session := p.GetSession()

	if userID, ok := session.Get("userID").(int64); ok && userID > 0 {
		user := goapp.LoadO[model.User](userID)

		if user == nil {
			//Session user not found, restart session
			session.Clear()
			session.Save()

			p.GetGinContext().Redirect(http.StatusSeeOther, "/")
			return false
		} else {
			p.Main("Already authenticated")
		}
	} else {
		fc := p.FormContext().
			SetParam("Session", session).
			SetRedirect(p.GetGinContext().DefaultQuery("destination", "/"))

		formOut := Forms.Login.Render(fc)

		if formOut.IsEmpty() {
			return false
		} else {
			p.Main(formOut)
		}
	}

	return true
}

func PageMyAccount(p *PageBuilderOLD) bool {
	p.Title(p.User().GetDisplayName())

	fc := p.FormContext().SetRedirect("/").
		SetParam("User", p.User())

	formOut := Forms.MyAccount.Render(fc)

	if formOut.IsEmpty() {
		return false
	} else {
		p.Main(formOut)
	}

	return true
}
