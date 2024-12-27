package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mtweb"
)

func PageAdminChecklists(p *PageBuilderOLD) bool {
	p.Main(
		mtweb.NewBtnPanel().Class("mb-3").AddIconBtn(
			"/", "home", "Home",
		).AddIconBtn(
			"/admin/checklists/0/edit", "plus", "Create checklist",
		),
	)

	cardList := mtweb.NewCardList()

	list := goapp.LoadOL[model.Checklist]()

	for _, cl := range list {
		card := mtweb.NewCard().
			Header(
				dhtmlbs.NewJustifiedLR().
					L(cl.Name).
					R(
						mtweb.NewEditBtn(fmt.Sprintf("/admin/checklists/%d/edit", cl.ID)),
					).
					R(
						mtweb.NewDeleteBtn(fmt.Sprintf("/admin/checklists/%d/delete", cl.ID), ""),
					),
			).
			Body("body")

		cardList.Add(card)
	}

	p.Main(cardList)
	return true
}

func PageAdminChecklistEdit(p *PageBuilderOLD) bool {
	cl := model.LoadChecklist(p.GetGinContext().Param("id"))

	if cl == nil {
		cl = &model.Checklist{}
		p.Title("New checklist")
	} else {
		p.Title("Edit checklist: " + cl.Name)
	}

	fc := p.FormContext().SetRedirect("/admin/checklists").
		SetParam("Checklist", cl)

	formOut := Forms.AdminChecklist.Render(fc)

	if formOut.IsEmpty() {
		return false
	} else {
		p.Main(formOut)
	}

	return true
}

func webAdminChecklistDelete(c *gin.Context) {
	if cl := model.LoadChecklist(c.Param("id")); cl != nil {
		goapp.DeleteObject(cl)
	}

	c.Redirect(http.StatusFound, "/admin/checklists")
}
