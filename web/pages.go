package web

import (
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/mtweb"
)

func PageDashboard(p *PageTemplate) {
	cards_list := dhtml.Div().Class("row row-cols-md-2 g-3")

	cards_list.Append(dhtml.Div().Class("col").Append(
		mtweb.NewCard().Header(mtweb.Icon("vial").Label("Experiment")).
			Body(dhtml.Div().Text("Html renderer ").Append(dhtml.NewLink("/experiment").Label("experiment")).Text(" link.")).
			Body(dhtml.Div().Text("Forms ").Append(dhtml.NewLink("/form").Label("experiment")).Text(" link.")),
	))

	cards_list.Append(dhtml.Div().Class("col").Append(
		mtweb.NewCard().Header(mtweb.Icon("list-check").Label("Active checklists")).
			Body("Some content"),
	))

	cards_list.Append(dhtml.Div().Class("col").Append(
		mtweb.NewCard().Header(mtweb.Icon("user-check").Label("My issues")).
			Body("Some content"),
	))

	cards_list.Append(dhtml.Div().Class("col").Append(
		mtweb.NewCard().Header(mtweb.Icon("chart-simple").Label("Statistics")).
			Body("Some content"),
	))

	cards_list.Append(dhtml.Div().Class("col").Append(
		mtweb.NewCard().Header(mtweb.Icon("cog").Label("System management")).
			Body(
				dhtml.Div().Append(
					dhtml.NewLink("/admin/checklists").Label(mtweb.Icon("list-check").Label("Checklists")),
				).Append(" (administration)"),
			),
	))
	p.Main(cards_list)
}
