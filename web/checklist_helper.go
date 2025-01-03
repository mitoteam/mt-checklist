package web

import (
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/mt-checklist/model"
)

func renderChecklistItemBody(item *model.ChecklistItem) (out dhtml.HtmlPiece) {
	if item.Body == "" {
		out.Append(dhtml.Div().Append(dhtml.EmptyLabel("no description")))
	} else {
		out.Append(dhtml.Div().Append(MdEngine.ToDhtml(item.Body)))
	}

	out.Append(dhtml.RenderValue("Responsible", model.LoadUser(item.ResponsibleID).GetDisplayName()).Class("mt-3"))

	return out
}
