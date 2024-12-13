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
		mtweb.NewBtnPanel().Class("mb-3").AddIconBtn(
			"/", "home", "Home",
		).AddIconBtn(
			"/admin/checklists/0/edit", "plus", "Create checklist",
		),
	)

	cardList := mtweb.NewCardList()

	list := goappbase.LoadOL[model.Checklist]()

	for _, cl := range list {
		card := mtweb.NewCard().
			Header(
				mtweb.NewJustifiedLR().
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

// ====================== user management ===================
func PageAdminUsers(p *PageBuilder) bool {
	p.Main(
		mtweb.NewBtnPanel().Class("mb-3").AddIconBtn(
			"/", "home", "Home",
		).AddIconBtn(
			"/admin/users/0/edit", "plus", "Create user",
		),
	)

	return true
}

func PageAdminUserEdit(p *PageBuilder) bool {
	user := goappbase.LoadOrCreateO[model.User](p.GetGinContext().Param("id"))

	if user == nil {
		p.Title("New user")
	} else {
		p.Title("Edit user: " + user.DisplayName)
	}

	fc := p.FormContext().SetRedirect("/admin/users").
		SetParam("User", user)

	formOut := dhtml.FormManager.RenderForm("admin_user_edit", fc)

	if formOut.IsEmpty() {
		return false
	} else {
		p.Main(formOut)
	}

	return true
}
