package web

import (
	"net/http"
	"time"

	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mttools"
	"github.com/mitoteam/mtweb"
)

type AdminController struct {
	mbr.ControllerBase
}

var AdminCtl *AdminController

func init() {
	AdminCtl = &AdminController{}
	AdminCtl.With(AuthMiddleware)
	AdminCtl.With(AdminRoleMiddleware)
}

func (c *AdminController) renderToolbar() *mtweb.BtnPanelElement {
	return mtweb.NewBtnPanel().Class("mb-3").AddIconBtn(
		mbr.Url(RootCtl.Home), iconHome, "Home",
	)
}

// ====================== user management ===================

func (c *AdminController) Users() mbr.Route {
	return mbr.Route{
		PathPattern: "/users",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			p.Main(
				c.renderToolbar().
					AddIconBtn(
						mbr.Url(AdminCtl.UserEdit, "user_id", 0), "plus", "Create user",
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

				actions.Append(mtweb.NewEditBtn(mbr.Url(AdminCtl.UserEdit, "user_id", user.ID)))
				actions.Append(
					mtweb.NewBtn().Class("btn-sm p-1").Href(mbr.Url(AdminCtl.UserPassword, "user_id", user.ID)).
						Title("Change password").Label(mtweb.Icon("key")),
				)
				actions.Append(mtweb.NewDeleteBtn(mbr.Url(AdminCtl.UserDelete, "user_id", user.ID), ""))

				row.Cell(actions)
			}

			return nil
		}),
	}
}

func (c *AdminController) UserEdit() mbr.Route {
	return mbr.Route{
		PathPattern: "/users/{user_id}/edit",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			user := goapp.LoadOrCreateO[model.User](p.ctx.Request().PathValue("user_id"))

			if user == nil {
				p.Title("New user")
			} else {
				p.Title("Edit user: " + user.DisplayName)
			}

			fc := p.FormContext().SetRedirect(mbr.Url(AdminCtl.Users)).
				SetArg("User", user)

			p.Main(formAdminUserEdit.Render(fc))
			return nil
		}),
	}
}

func (c *AdminController) UserPassword() mbr.Route {
	return mbr.Route{
		PathPattern: "/users/{user_id}/password",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			user := goapp.LoadOMust[model.User](p.ctx.Request().PathValue("user_id"))

			p.Title("User password: " + user.DisplayName)

			fc := p.FormContext().SetRedirect(mbr.Url(AdminCtl.Users)).
				SetArg("User", user)

			p.Main(formAdminUserPassword.Render(fc))

			return nil
		}),
	}
}

func (c *AdminController) UserDelete() mbr.Route {
	return mbr.Route{
		PathPattern: "/users/{user_id}/delete",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			goapp.DeleteObject(goapp.LoadOMust[model.User](p.ctx.Request().PathValue("user_id")))
			p.ctx.RedirectRoute(http.StatusFound, AdminCtl.Users)

			return nil
		}),
	}
}

// ====================== checklist templates management ===================

func (c *AdminController) renderTemplatesToolbar() *mtweb.BtnPanelElement {
	return c.renderToolbar().
		AddIconBtn(
			mbr.Url(AdminCtl.Templates), iconTemplate, "All Templates",
		)
}

func (c *AdminController) Templates() mbr.Route {
	return mbr.Route{
		PathPattern: "/template",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			p.Main(
				c.renderToolbar().
					AddIconBtn(
						mbr.Url(AdminCtl.TemplateEdit, "template_id", 0), "plus", "Create template",
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
						mtweb.NewBtn().Href(mbr.Url(AdminCtl.TemplateItemList, "template_id", t.ID)).Class("btn-sm").
							Label(mtweb.Icon("list-check").Label(t.ItemCount())),
					)

				var actions dhtml.HtmlPiece

				actions.Append(mtweb.NewEditBtn(mbr.Url(AdminCtl.TemplateEdit, "template_id", t.ID))).
					Append(mtweb.NewDeleteBtn(mbr.Url(AdminCtl.TemplateDelete, "template_id", t.ID), ""))

				row.Cell(actions)
			}

			p.Main(table)

			return nil
		}),
	}
}

func (c *AdminController) TemplateEdit() mbr.Route {
	return mbr.Route{
		PathPattern: "/template/{template_id}/edit",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			t := goapp.LoadOrCreateO[model.ChecklistTemplate](p.ctx.Request().PathValue("template_id"))

			if t == nil {
				p.Title("New template")
			} else {
				p.Title("Edit template: " + t.Name)
			}

			fc := p.FormContext().SetRedirect(mbr.Url(AdminCtl.Templates)).SetArg("Template", t)

			p.Main(formAdminChecklistTemplate.Render(fc))

			return nil
		}),
	}
}

func (c *AdminController) TemplateDelete() mbr.Route {
	return mbr.Route{
		PathPattern: "/template/{template_id}/delete",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			goapp.DeleteObject(goapp.LoadOMust[model.ChecklistTemplate](p.ctx.Request().PathValue("template_id")))
			p.ctx.RedirectRoute(http.StatusFound, AdminCtl.Templates)
			return nil
		}),
	}
}

func (c *AdminController) TemplateItemList() mbr.Route {
	return mbr.Route{
		PathPattern: "/template/{template_id}/items",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			t := goapp.LoadOrCreateO[model.ChecklistTemplate](p.ctx.Request().PathValue("template_id"))

			p.Main(
				c.renderTemplatesToolbar().
					AddIconBtn(
						mbr.Url(AdminCtl.TemplateItemEdit, "template_id", t.ID, "item_id", 0), "plus", "Add item",
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

				actions.Append(mtweb.NewEditBtn(mbr.Url(AdminCtl.TemplateItemEdit, "template_id", t.ID, "item_id", item.ID)))
				actions.Append(mtweb.NewDeleteBtn(mbr.Url(AdminCtl.TemplateItemDelete, "template_id", t.ID, "item_id", item.ID), ""))

				row.Cell(actions)
			}

			p.Main(table)

			return nil
		}),
	}
}

func (c *AdminController) TemplateItemEdit() mbr.Route {
	return mbr.Route{
		PathPattern: "/template/{template_id}/item/{item_id}/edit",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			t := goapp.LoadOrCreateO[model.ChecklistTemplate](p.ctx.Request().PathValue("template_id"))
			item := goapp.LoadOrCreateO[model.ChecklistTemplateItem](p.ctx.Request().PathValue("item_id"))

			if item.ID == 0 {
				item.ChecklistTemplateID = t.ID
				p.Title("New item")
			} else {
				mttools.AssertEqual(item.ChecklistTemplateID, t.ID)
				p.Title("Edit item: " + item.Caption)
			}

			fc := p.FormContext().SetRedirect(mbr.Url(AdminCtl.TemplateItemList, "template_id", t.ID)).
				SetArg("Item", item)

			p.Main(formAdminChecklistTemplateItem.Render(fc))

			return nil
		}),
	}
}

func (c *AdminController) TemplateItemDelete() mbr.Route {
	return mbr.Route{
		PathPattern: "/template/{template_id}/item/{item_id}/delete",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			t := goapp.LoadOrCreateO[model.ChecklistTemplate](p.ctx.Request().PathValue("template_id"))
			item := goapp.LoadOrCreateO[model.ChecklistTemplateItem](p.ctx.Request().PathValue("item_id"))

			mttools.AssertEqual(item.ChecklistTemplateID, t.ID)
			goapp.DeleteObject(item)

			p.ctx.RedirectRoute(http.StatusFound, AdminCtl.TemplateItemList, "template_id", t.ID)

			return nil
		}),
	}
}

// ====================== checklists admin management ===================

func (c *AdminController) Checklists() mbr.Route {
	return mbr.Route{
		PathPattern: "/checklists",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			p.Main(
				c.renderToolbar().
					AddIconBtn(
						mbr.Url(AdminCtl.ChecklistEdit, "checklist_id", 0), "plus", "Create checklist",
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
								mtweb.NewEditBtn(mbr.Url(AdminCtl.ChecklistEdit, "checklist_id", cl.ID)),
							).
							R(
								mtweb.NewDeleteBtn(mbr.Url(AdminCtl.ChecklistDelete, "checklist_id", cl.ID), ""),
							),
					).
					Body("body")

				cardList.Add(card)
			}

			p.Main(cardList)

			return nil
		}),
	}
}

func (c *AdminController) ChecklistEdit() mbr.Route {
	return mbr.Route{
		PathPattern: "/checklist/{checklist_id}/edit",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			cl := goapp.LoadO[model.Checklist](p.ctx.Request().PathValue("checklist_id"))

			if cl == nil {
				cl = &model.Checklist{}
				p.Title("New checklist")
			} else {
				p.Title("Edit checklist: " + cl.Name)
			}

			fc := p.FormContext().SetRedirect(mbr.Url(AdminCtl.Checklists)).SetArg("Checklist", cl)
			p.Main(formAdminChecklist.Render(fc))

			return nil
		}),
	}
}

func (c *AdminController) ChecklistDelete() mbr.Route {
	return mbr.Route{
		PathPattern: "/checklist/{checklist_id}/delete",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			goapp.DeleteObject(model.LoadChecklist(p.ctx.Request().PathValue("checklist_id")))
			p.ctx.RedirectRoute(http.StatusFound, AdminCtl.Checklists)

			return nil
		}),
	}
}
