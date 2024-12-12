package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/goappbase"
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
					).
					R(
						dhtml.NewConfirmLink(fmt.Sprintf("/admin/checklists/%d/delete", cl.ID), "").
							Label(mtweb.Icon("trash").Class("text-danger")),
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

func webAdminChecklistDelete(c *gin.Context) {
	if cl := model.LoadChecklist(c.Param("id")); cl != nil {
		goappbase.DeleteObject(cl)
	}

	c.Redirect(http.StatusFound, "/admin/checklists")
}
