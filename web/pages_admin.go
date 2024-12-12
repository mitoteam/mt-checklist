package web

import (
	"fmt"

	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mtweb"
)

func PageAdminChecklists(p *PageBuilder) bool {
	p.Main(
		dhtml.Div().Class("mb-3 p-3 border").
			Append(
				dhtml.NewLink("/admin/checklists/0/edit").Label(mtweb.Icon("plus").Label("Create checklist")),
			),
	)

	cardList := mtweb.NewCardList()

	for _, cl := range model.GetChecklistsList() {
		card := mtweb.NewCard().
			Header(
				mtweb.NewJustifiedLR().
					L(cl.Name).
					R(
						dhtml.NewLink(fmt.Sprintf("/admin/checklists/%d/edit", cl.ID)).Label(mtweb.Icon("edit")),
					),
			).
			Body("body")

		cardList.Add(card)
	}

	p.Main(cardList)
	return true
}

func PageAdminChecklistEdit(p *PageBuilder) bool {
	cl := model.LoadChecklist(p.GetGinContext().Param("id"))

	if cl == nil {
		cl = &model.Checklist{}
		p.Title("New checklist")
	} else {
		p.Title("Edit checklist: " + cl.Name)
	}

	fc := p.FormContext().SetRedirect("/admin/checklists").
		SetParam("Checklist", cl)

	formOut := dhtml.FormManager.RenderForm("admin_checklist_edit", fc)

	if formOut.IsEmpty() {
		return false
	} else {
		p.Main(formOut)
	}

	return true
}
