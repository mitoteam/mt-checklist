package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-checklist/app"
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

// base path fo all admin routes
func (root *RootController) AdminSubroutes() mbr.Route {
	return mbr.Route{PathPattern: "/admin", ChildController: AdminCtl}
}

// toolbar helper
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
			p.Title("Users").Main(
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

			for _, user := range goapp.LoadOL[model.User]() {
				row := table.NewRow()

				row.Cell(user.UserName)
				row.Cell(user.DisplayName)
				row.Cell(mtweb.IconYesNo(user.IsActive))
				row.Cell(mtweb.IconYesNo(user.IsAdmin()))

				if user.LastLogin != nil {
					row.Cell(user.LastLogin.Format(time.DateTime))
				} else {
					row.Cell(mtweb.IconNo())
				}

				var actions dhtml.HtmlPiece

				actions.Append(mtweb.NewEditBtn(mbr.Url(AdminCtl.UserEdit, "user_id", user.ID)))
				actions.Append(
					dhtmlbs.NewBtn().Class("btn-sm p-1").Href(mbr.Url(AdminCtl.UserPassword, "user_id", user.ID)).
						Title("Change password").Label(mtweb.Icon("key")),
				)
				actions.Append(mtweb.NewDeleteBtn(mbr.Url(AdminCtl.UserDelete, "user_id", user.ID), ""))

				row.Cell(actions)
			}

			p.Main(mtweb.RenderTableCount(table, "Users count"))

			p.Main(table)

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
			p.RedirectRoute(AdminCtl.Users)

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
			p.Title("Checklist templates").Main(
				c.renderToolbar().
					AddIconBtn(
						mbr.Url(AdminCtl.TemplateEdit, "template_id", 0), "plus", "Create template",
					),
			)

			table := mtweb.NewTable().EmptyLabel("no checklist templates created yet").
				Header("Template Name").
				Header("Checklist Name").
				Header("Items").
				Header("") //actions

			for _, t := range goapp.LoadOL[model.Template]() {
				row := table.NewRow()

				row.Cell(t.Name)
				row.Cell(t.ChecklistName)
				row.Cell(
					dhtmlbs.NewBtn().Href(mbr.Url(AdminCtl.TemplateItemsList, "template_id", t.ID)).Class("btn-sm").
						Label(mtweb.Icon("list-check").Label(t.ItemCount())),
				)

				var actions dhtml.HtmlPiece

				actions.
					Append(mtweb.NewEditBtn(mbr.Url(AdminCtl.TemplateEdit, "template_id", t.ID))).
					Append(mtweb.NewDeleteBtn(mbr.Url(AdminCtl.TemplateDelete, "template_id", t.ID), "")).
					Append(
						mtweb.NewIconBtn(mbr.Url(AdminCtl.TemplateCreateChecklist, "template_id", t.ID), "plus", "Create checklist").
							Class("btn-success btn-sm").
							Confirm(fmt.Sprintf("Do you want to create new checklist from %s template?", t.Name)),
					)

				row.Cell(actions)
			}

			p.Main(mtweb.RenderTableCount(table, "Templates count"))
			p.Main(table)

			return nil
		}),
	}
}

func (c *AdminController) TemplateEdit() mbr.Route {
	return mbr.Route{
		PathPattern: "/template/{template_id}/edit",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			t := goapp.LoadOrCreateO[model.Template](p.ctx.Request().PathValue("template_id"))

			if t == nil {
				p.Title("New template")
			} else {
				p.Title("Edit template: " + t.Name)
			}

			fc := p.FormContext().SetRedirect(mbr.Url(AdminCtl.Templates)).SetArg("Template", t)

			p.Main(formAdminTemplate.Render(fc))

			return nil
		}),
	}
}

func (c *AdminController) TemplateRenumber() mbr.Route {
	return mbr.Route{
		PathPattern: "/template/{template_id}/renumber",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			t := goapp.LoadOrCreateO[model.Template](p.ctx.Request().PathValue("template_id"))

			p.Title("Renumber template items: " + t.Name)

			fc := p.FormContext().SetArg("Template", t).
				SetRedirect(mbr.Url(AdminCtl.TemplateItemsList, "template_id", t.ID))

			p.Main(formAdminTemplateRenumber.Render(fc))

			return nil
		}),
	}
}

func (c *AdminController) TemplateDelete() mbr.Route {
	return mbr.Route{
		PathPattern: "/template/{template_id}/delete",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			goapp.DeleteObject(goapp.LoadOMust[model.Template](p.ctx.Request().PathValue("template_id")))
			p.RedirectRoute(AdminCtl.Templates)
			return nil
		}),
	}
}

func (c *AdminController) TemplateItemsList() mbr.Route {
	return mbr.Route{
		PathPattern: "/template/{template_id}/items",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			t := goapp.LoadOrCreateO[model.Template](p.ctx.Request().PathValue("template_id"))

			p.Title("Checklist template items").Main(
				c.renderTemplatesToolbar().
					AddIconBtn(
						mbr.Url(AdminCtl.TemplateItemEdit, "template_id", t.ID, "item_id", 0), "plus", "Add item",
					).
					AddIconBtn(
						mbr.Url(AdminCtl.TemplateRenumber, "template_id", t.ID), "list-ol", "Renumber items",
					),
			).Main(
				dhtml.RenderValue(
					"Template",
					dhtml.NewLink(mbr.Url(AdminCtl.TemplateEdit, "template_id", t.ID)).Label(mtweb.Icon(iconTemplate).Label(t.Name)),
				).Class("mb-3"),
			)

			table := dhtml.NewTable().Class("table table-hover table-sm").EmptyLabel("no items added yet").
				Header("Caption / Body").
				Header("Responsible").
				Header("Depends").
				Header("Sort Order").
				Header("Weight").
				Header("") //actions

			for _, item := range t.Items() {
				row := table.NewRow()

				//caption and body
				cellOut := dhtml.Piece(dhtml.Div().Class("fw-bold mb-1").Append(item.Caption))
				cellOut.Append(dhtml.Div().Class("small text-prewrap").Append(item.Body))
				row.Cell(cellOut)

				//responsible
				row.Cell(item.GetResponsible().GetDisplayName())

				//dependencies
				flexContainer := dhtml.Div().Class("d-flex")
				if item.DependenciesCount() > 0 {
					depsList := dhtml.NewUnorderedList()

					for _, dep := range item.DependenciesList() {
						depsList.AppendItem(dhtml.NewListItem().Append(dep.GetRequireTemplateItem().Caption))
					}

					flexContainer.Append(depsList)
				} else {
					flexContainer.Append(dhtml.Div().Append(dhtml.EmptyLabel("no dependencies")))
				}
				flexContainer.Append(dhtml.Div().Class("ms-1").Append(
					mtweb.NewIconBtn(
						mbr.Url(AdminCtl.TemplateItemDependencies, "template_id", t.ID, "item_id", item.ID),
						iconDependencies, "",
					).Class("btn-sm p-1").Title("Edit dependencies"),
				))
				row.Cell(flexContainer)

				row.Cell(item.SortOrder)
				row.Cell(item.Weight)

				var actions dhtml.HtmlPiece

				actions.Append(mtweb.NewEditBtn(mbr.Url(AdminCtl.TemplateItemEdit, "template_id", t.ID, "item_id", item.ID)))
				actions.Append(mtweb.NewDeleteBtn(mbr.Url(AdminCtl.TemplateItemDelete, "template_id", t.ID, "item_id", item.ID), ""))

				row.Cell(actions)
			}

			p.Main(mtweb.RenderTableCount(table, "Template items count"))
			p.Main(table)

			return nil
		}),
	}
}

func (c *AdminController) TemplateItemEdit() mbr.Route {
	return mbr.Route{
		PathPattern: "/template/{template_id}/item/{item_id}/edit",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			t := goapp.LoadOrCreateO[model.Template](p.ctx.Request().PathValue("template_id"))
			item := goapp.LoadOrCreateO[model.TemplateItem](p.ctx.Request().PathValue("item_id"))

			if item.ID == 0 {
				p.Title("New item")

				item.TemplateID = t.ID
				item.ResponsibleID = p.User().ID //current user by default
				item.Weight = 1
				item.SortOrder = t.MaxItemSortOrder() + app.App.AppSettings.(*app.AppSettingsType).SortOrderStep
			} else {
				mttools.AssertEqual(item.TemplateID, t.ID)

				p.Title("Edit item: " + item.Caption)
			}

			fc := p.FormContext().SetRedirect(mbr.Url(AdminCtl.TemplateItemsList, "template_id", t.ID)).
				SetArg("Item", item)

			p.Main(formAdminChecklistTemplateItem.Render(fc))

			return nil
		}),
	}
}

func (c *AdminController) TemplateItemDependencies() mbr.Route {
	return mbr.Route{
		PathPattern: "/template/{template_id}/item/{item_id}/deps",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			t := goapp.LoadOMust[model.Template](p.ctx.Request().PathValue("template_id"))
			item := goapp.LoadOMust[model.TemplateItem](p.ctx.Request().PathValue("item_id"))

			mttools.AssertEqual(item.TemplateID, t.ID)
			p.Title("Item dependencies: " + item.Caption)

			fc := p.FormContext().SetRedirect(mbr.Url(AdminCtl.TemplateItemsList, "template_id", t.ID)).
				SetArg("Item", item)

			p.Main(formAdminChecklistTemplateItemDeps.Render(fc))

			return nil
		}),
	}
}

func (c *AdminController) TemplateItemDelete() mbr.Route {
	return mbr.Route{
		PathPattern: "/template/{template_id}/item/{item_id}/delete",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			t := goapp.LoadOrCreateO[model.Template](p.ctx.Request().PathValue("template_id"))
			item := goapp.LoadOrCreateO[model.TemplateItem](p.ctx.Request().PathValue("item_id"))

			mttools.AssertEqual(item.TemplateID, t.ID)
			goapp.DeleteObject(item)

			p.RedirectRoute(AdminCtl.TemplateItemsList, "template_id", t.ID)

			return nil
		}),
	}
}

func (c *AdminController) TemplateCreateChecklist() mbr.Route {
	return mbr.Route{
		PathPattern: "/template/{template_id}/create-checklist",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			template := goapp.LoadOrCreateO[model.Template](p.ctx.Request().PathValue("template_id"))

			checklist := createChecklistFromTemplate(template, p.User())

			p.RedirectRoute(AdminCtl.ChecklistEdit, "checklist_id", checklist.ID)

			return nil
		}),
	}
}

// ====================== checklists admin management ===================

func (c *AdminController) renderChecklistsToolbar() *mtweb.BtnPanelElement {
	return c.renderToolbar().
		AddIconBtn(
			mbr.Url(AdminCtl.Checklists), iconChecklist, "All Checklists",
		)
}

func (c *AdminController) Checklists() mbr.Route {
	return mbr.Route{
		PathPattern: "/checklists",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			p.Title("Checklists (admin)").Main(
				c.renderToolbar().
					AddIconBtn(
						mbr.Url(AdminCtl.ChecklistEdit, "checklist_id", 0), "plus", "Create checklist",
					),
			)

			table := mtweb.NewTable().
				Header("Active").
				Header("Name").
				Header("Description").
				Header("Created").
				Header("Items").
				Header("") //actions

			list := goapp.LoadOL[model.Checklist]()

			for _, cl := range list {
				row := table.NewRow()

				row.Cell(mtweb.IconYesNo(cl.IsActive()))
				row.Cell(cl.Name)
				row.Cell(cl.Description).Class("small text-muted")

				cellOut := dhtml.Piece(cl.GetCreatedBy().GetDisplayName())
				cellOut.Append(dhtml.Div().Append(mtweb.RenderTimestamp(cl.CreatedAt)))
				row.Cell(cellOut)

				row.Cell(
					mtweb.NewIconBtn(mbr.Url(AdminCtl.ChecklistItemsList, "checklist_id", cl.ID), iconChecklist, cl.ItemCount()).Class("btn-sm"),
				)

				var actions dhtml.HtmlPiece

				actions.Append(
					mtweb.NewIconBtn(mbr.Url(ChecklistCtl.Checklist, "checklist_id", cl.ID), iconView, nil).
						Class("btn-sm px-1").Title("View checklist"),
				)
				actions.Append(mtweb.NewEditBtn(mbr.Url(AdminCtl.ChecklistEdit, "checklist_id", cl.ID)))
				actions.Append(mtweb.NewDeleteBtn(mbr.Url(AdminCtl.ChecklistDelete, "checklist_id", cl.ID), ""))

				row.Cell(actions)
			}

			p.Main(mtweb.RenderTableCount(table, "Checklists count"))
			p.Main(table)

			return nil
		}),
	}
}

func (c *AdminController) ChecklistEdit() mbr.Route {
	return mbr.Route{
		PathPattern: "/checklist/{checklist_id}/edit",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			cl := goapp.LoadOrCreateO[model.Checklist](p.ctx.Request().PathValue("checklist_id"))

			if cl.ID == 0 {
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

func (c *AdminController) ChecklistItemsList() mbr.Route {
	return mbr.Route{
		PathPattern: "/checklist/{checklist_id}/items",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			cl := model.LoadChecklist(p.ctx.Request().PathValue("checklist_id"))

			p.Title("Checklist Items: " + cl.Name).Main(
				c.renderChecklistsToolbar().
					AddIconBtn(
						mbr.Url(AdminCtl.ChecklistItemEdit, "checklist_id", cl.ID, "item_id", 0), "plus", "Add item",
					),
			).Main(
				dhtml.RenderValue(mtweb.Icon(iconChecklist).Label("Checklist"), cl.Name).Class("mb-3"),
			)

			table := mtweb.NewTable().
				Header("Caption / Body").
				Header("Responsible").
				Header("Done").
				Header("Depends").
				Header("Sort Order").
				Header("Weight").
				Header("") //actions

			for _, item := range cl.Items() {
				row := table.NewRow()

				//caption and body
				cellOut := dhtml.Piece(dhtml.Div().Class("fw-bold mb-1").Append(item.Caption))
				cellOut.Append(dhtml.Div().Class("small text-prewrap").Append(item.Body))
				row.Cell(cellOut)

				//responsible
				row.Cell(item.GetResponsible().GetDisplayName())

				//done
				cellOut.Clear()
				if item.DoneAt == nil {
					cellOut.Append(mtweb.IconNo())
				} else {
					cellOut.Append(item.GetDoneBy().GetDisplayName()).Append(
						dhtml.Div().Append(mtweb.RenderTimestamp(*item.DoneAt)),
					)
				}
				row.Cell(cellOut)

				//dependencies
				flexC := dhtml.Div().Class("d-flex")
				if item.DependenciesCount() > 0 {
					depsList := dhtml.NewUnorderedList()

					for _, dep := range item.DependenciesList() {
						depsList.AppendItem(dhtml.NewListItem().Append(dep.GetRequireChecklistItem().Caption))
					}

					flexC.Append(depsList)
				} else {
					flexC.Append(dhtml.Div().Append(dhtml.EmptyLabel("no dependencies")))
				}
				flexC.Append(dhtml.Div().Class("ms-2").Append(
					mtweb.NewIconBtn(
						mbr.Url(AdminCtl.ChecklistItemDependencies, "checklist_id", cl.ID, "item_id", item.ID), iconDependencies, "",
					).Class("btn-sm p-1").Title("Edit dependencies"),
				))
				row.Cell(flexC)

				row.Cell(item.SortOrder)
				row.Cell(item.Weight)

				var actions dhtml.HtmlPiece

				actions.Append(mtweb.NewEditBtn(mbr.Url(AdminCtl.ChecklistItemEdit, "checklist_id", cl.ID, "item_id", item.ID)))
				actions.Append(mtweb.NewDeleteBtn(mbr.Url(AdminCtl.ChecklistItemDelete, "checklist_id", cl.ID, "item_id", item.ID), ""))

				row.Cell(actions)
			}

			p.Main(mtweb.RenderTableCount(table, "Checklist items count"))
			p.Main(table)

			return nil
		}),
	}
}

func (c *AdminController) ChecklistItemEdit() mbr.Route {
	return mbr.Route{
		PathPattern: "/checklist/{checklist_id}/item/{item_id}/edit",
		Method:      "GET POST",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			cl := model.LoadChecklist(p.ctx.Request().PathValue("checklist_id"))
			item := goapp.LoadOrCreateO[model.ChecklistItem](p.ctx.Request().PathValue("item_id"))

			if item.ID == 0 {
				//new item
				item.ChecklistID = cl.ID
				item.ResponsibleID = p.User().ID //current user by default
				item.Weight = 1
				item.SortOrder = cl.MaxItemSortOrder() + app.App.AppSettings.(*app.AppSettingsType).SortOrderStep
			} else {
				//existing
				mttools.AssertEqual(item.ChecklistID, cl.ID)
			}

			fc := p.FormContext().SetArg("Item", item).SetRedirect(mbr.Url(AdminCtl.ChecklistItemsList, "checklist_id", cl.ID))
			p.Main(formAdminChecklistItem.Render(fc))

			return nil
		}),
	}
}

func (c *AdminController) ChecklistItemDelete() mbr.Route {
	return mbr.Route{
		PathPattern: "/checklist/{checklist_id}/item/{item_id}/delete",
		Method:      "GET",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			cl := model.LoadChecklist(p.ctx.Request().PathValue("checklist_id"))
			item := goapp.LoadOrCreateO[model.ChecklistItem](p.ctx.Request().PathValue("item_id"))

			mttools.AssertEqual(item.ChecklistID, cl.ID)
			goapp.DeleteObject(item)

			p.ctx.RedirectRoute(http.StatusFound, AdminCtl.ChecklistItemsList, "checklist_id", cl.ID)

			return nil
		}),
	}
}

func (c *AdminController) ChecklistItemDependencies() mbr.Route {
	return mbr.Route{
		PathPattern: "/checklist/{checklist_id}/item/{item_id}/deps",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			cl := goapp.LoadOMust[model.Checklist](p.ctx.Request().PathValue("checklist_id"))
			item := goapp.LoadOMust[model.ChecklistItem](p.ctx.Request().PathValue("item_id"))

			mttools.AssertEqual(item.ChecklistID, cl.ID)
			p.Title("Item dependencies: " + item.Caption)

			fc := p.FormContext().SetRedirect(mbr.Url(AdminCtl.ChecklistItemsList, "checklist_id", cl.ID)).
				SetArg("Item", item)

			p.Main(formAdminChecklistItemDeps.Render(fc))

			return nil
		}),
	}
}
