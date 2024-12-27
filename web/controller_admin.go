package web

import (
	"net/http"
	"time"

	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-checklist/model"
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

// ====================== user management ===================

func (c *AdminController) Users() mbr.Route {
	return mbr.Route{
		PathPattern: "/users",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			p.Main(
				mtweb.NewBtnPanel().Class("mb-3").AddIconBtn(
					mbr.Url(RootCtl.Home), iconHome, "Home",
				).AddIconBtn(
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

			return true
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
			p.ctx.Redirect(http.StatusFound, mbr.Url(AdminCtl.Users))

			return nil
		}),
	}
}
