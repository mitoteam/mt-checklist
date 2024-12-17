package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/goapp"
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

	list := goapp.LoadOL[model.Checklist]()

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

	formOut := dhtml.FormManager.RenderForm(Forms.AdminChecklist, fc)

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

// ====================== user management ===================

func PageAdminUsers(p *PageBuilder) bool {
	p.Main(
		mtweb.NewBtnPanel().Class("mb-3").AddIconBtn(
			"/", "home", "Home",
		).AddIconBtn(
			"/admin/users/0/edit", "plus", "Create user",
		),
	)

	table := dhtml.NewTable().Class("table table-hover table-sm").
		Header("Username").
		Header("Display name").
		Header("Active").
		Header("Admin").
		Header("Last Login").
		Header("")

	p.Main(table)

	for _, user := range goapp.LoadOL[model.User]() {
		row := table.NewRow()

		row.Cell(user.UserName).
			Cell(user.DisplayName).
			Cell(mtweb.IconYesNo(user.IsActive)).
			Cell(mtweb.IconYesNo(user.IsAdmin()))

		if user.LastLogin != nil {
			row.Cell(user.LastLogin.Format(time.DateTime))
		} else {
			row.Cell(mtweb.IconNo())
		}

		var actions dhtml.HtmlPiece

		actions.Append(mtweb.NewEditBtn(fmt.Sprintf("/admin/users/%d/edit", user.ID)))
		actions.Append(
			mtweb.NewBtn().Class("btn-sm p-1").Href(fmt.Sprintf("/admin/users/%d/password", user.ID)).
				Title("Change password").Label(mtweb.Icon("key")),
		)
		actions.Append(mtweb.NewDeleteBtn(fmt.Sprintf("/admin/users/%d/delete", user.ID), ""))

		row.Cell(actions)
	}

	return true
}

func PageAdminUserEdit(p *PageBuilder) bool {
	user := goapp.LoadOrCreateO[model.User](p.GetGinContext().Param("id"))

	if user == nil {
		p.Title("New user")
	} else {
		p.Title("Edit user: " + user.DisplayName)
	}

	fc := p.FormContext().SetRedirect("/admin/users").
		SetParam("User", user)

	formOut := dhtml.FormManager.RenderForm(Forms.AdminUserEdit, fc)

	if formOut.IsEmpty() {
		return false
	} else {
		p.Main(formOut)
	}

	return true
}

func PageAdminUserPassword(p *PageBuilder) bool {
	user := goapp.LoadOMust[model.User](p.GetGinContext().Param("id"))
	p.Title("User password: " + user.DisplayName)

	fc := p.FormContext().SetRedirect("/admin/users").
		SetParam("User", user)

	formOut := dhtml.FormManager.RenderForm(Forms.AdminUserPassword, fc)

	if formOut.IsEmpty() {
		return false
	} else {
		p.Main(formOut)
	}

	return true
}

func webAdminUserDelete(c *gin.Context) {
	goapp.DeleteObject(goapp.LoadOMust[model.User](c.Param("id")))
	c.Redirect(http.StatusFound, "/admin/users")
}

// ====================== checklist templates management ===================

func PageAdminChecklistTemplates(p *PageBuilder) bool {
	p.Main(
		mtweb.NewBtnPanel().Class("mb-3").AddIconBtn(
			"/", "home", "Home",
		).AddIconBtn(
			"/admin/templates/0/edit", "plus", "Create template",
		),
	)

	table := dhtml.NewTable().Class("table table-hover table-sm").EmptyLabel("no checklist templates created yet").
		Header("Name").
		Header("Checklist Name").
		Header("Items").
		Header("") //actions

	for _, t := range goapp.LoadOL[model.ChecklistTemplate]() {
		row := table.NewRow()

		row.Cell(t.Name).
			Cell(t.ChecklistName).
			Cell(
				mtweb.NewBtn().Href(fmt.Sprintf("/admin/templates/%d/items", t.ID)).Class("btn-sm").
					Label(mtweb.Icon("list-check").Label(t.ItemCount())),
			)

		var actions dhtml.HtmlPiece

		actions.Append(mtweb.NewEditBtn(fmt.Sprintf("/admin/templates/%d/edit", t.ID)))
		actions.Append(mtweb.NewDeleteBtn(fmt.Sprintf("/admin/templates/%d/delete", t.ID), ""))

		row.Cell(actions)
	}

	p.Main(table)

	return true
}

func PageAdminChecklistTemplateEdit(p *PageBuilder) bool {
	t := goapp.LoadOrCreateO[model.ChecklistTemplate](p.GetGinContext().Param("id"))

	if t == nil {
		p.Title("New template")
	} else {
		p.Title("Edit template: " + t.Name)
	}

	fc := p.FormContext().SetRedirect("/admin/templates").
		SetParam("Template", t)

	formOut := dhtml.FormManager.RenderForm(Forms.AdminChecklistTemplate, fc)

	if formOut.IsEmpty() {
		return false
	} else {
		p.Main(formOut)
	}

	return true
}

func webAdminChecklistTemplateDelete(c *gin.Context) {
	goapp.DeleteObject(goapp.LoadOMust[model.ChecklistTemplate](c.Param("id")))
	c.Redirect(http.StatusFound, "/admin/templates")
}

func PageAdminChecklistTemplateItemList(p *PageBuilder) bool {
	t := goapp.LoadOMust[model.ChecklistTemplate](p.GetGinContext().Param("id"))

	p.Main(
		mtweb.NewBtnPanel().Class("mb-3").AddIconBtn(
			"/", "home", "Home",
		).AddIconBtn(
			"/admin/templates",
			iconTemplate, "All Templates",
		).AddIconBtn(
			fmt.Sprintf("/admin/templates/%d/items/0/edit", t.ID),
			"plus", "Add item",
		),
	).Main(
		dhtml.RenderValue("Template", t.Name).Class("mb-3"),
	)

	table := dhtml.NewTable().Class("table table-hover table-sm").EmptyLabel("no items added yet").
		Header("Caption").
		Header("Body").
		Header("Sort Order").
		Header("Weight").
		Header("") //actions

	for _, item := range t.Items() {
		row := table.NewRow()

		row.Cell(item.Caption).
			Cell(item.Body).
			Cell(item.SortOrder).
			Cell(item.Weight)

		var actions dhtml.HtmlPiece

		actions.Append(mtweb.NewEditBtn(fmt.Sprintf("/admin/templates/%d/items/%d/edit", t.ID, item.ID)))
		actions.Append(mtweb.NewDeleteBtn(fmt.Sprintf("/admin/templates/%d/items/%d/delete", t.ID, item.ID), ""))

		row.Cell(actions)
	}

	p.Main(table)

	return true
}

func PageAdminChecklistTemplateItemEdit(p *PageBuilder) bool {
	t := goapp.LoadOMust[model.ChecklistTemplate](p.GetGinContext().Param("id"))
	item := goapp.LoadOrCreateO[model.ChecklistTemplateItem](p.GetGinContext().Param("item_id"))

	if item.ID == 0 {
		p.Title("New item")
		item.ChecklistTemplateID = t.ID
	} else {
		if item.ChecklistTemplateID != t.ID {
			return false
		}

		p.Title("Edit item: " + item.Caption)
	}

	fc := p.FormContext().SetRedirect(fmt.Sprintf("/admin/templates/%d/items", t.ID)).
		SetParam("Item", item)

	formOut := dhtml.FormManager.RenderForm(Forms.AdminChecklistTemplateItem, fc)

	if formOut.IsEmpty() {
		return false
	} else {
		p.Main(formOut)
	}

	return true
}

func webAdminChecklistTemplateDeleteItem(c *gin.Context) {
	item := goapp.LoadOMust[model.ChecklistTemplateItem](c.Param("item_id"))
	t := goapp.LoadOMust[model.ChecklistTemplate](c.Param("id"))

	if item.ChecklistTemplateID != t.ID {
		goapp.DeleteObject(item)
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/admin/templates/%d/items", t.ID))
}
