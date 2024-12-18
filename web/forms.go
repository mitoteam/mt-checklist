package web

import (
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlform"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mtweb"
)

// helper type to have all form handlers in one place
type formsType struct {
	Login     *dhtmlform.FormHandler
	MyAccount *dhtmlform.FormHandler

	AdminUserEdit     *dhtmlform.FormHandler
	AdminUserPassword *dhtmlform.FormHandler

	AdminChecklist *dhtmlform.FormHandler

	AdminChecklistTemplate     *dhtmlform.FormHandler
	AdminChecklistTemplateItem *dhtmlform.FormHandler
}

var Forms formsType

func init() {
	Forms.Login = &dhtmlform.FormHandler{
		RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
			formBody.Append(dhtml.Div().Class("border bg-light p-3").Append(
				dhtmlform.NewTextInput("username").Placeholder("Username").Require(),
				dhtmlform.NewPasswordInput("password").Placeholder("Password").Require(),
				dhtmlform.NewSubmitBtn().Label(mtweb.Icon("arrow-right-to-bracket").Label("Sign In")),
			))
		},
		ValidateF: func(fd *dhtmlform.FormData) {
			if !fd.HasError() {
				username := fd.GetValue("username").(string)
				password := fd.GetValue("password").(string)

				user := model.AuthorizeUser(username, password)

				if user != nil {
					fd.SetParam("userID", user.ID)
				} else {
					if session, ok := fd.GetParam("Session").(sessions.Session); ok {
						session.Delete("userID")
						session.Save()
					}

					fd.SetError("", "User not found or wrong password given")
				}
			}
		},
		SubmitF: func(fd *dhtmlform.FormData) {
			if session, ok := fd.GetParam("Session").(sessions.Session); ok {
				session.Set("userID", fd.GetParam("userID").(int64))
				session.Save()
			}
		},
	}

	Forms.MyAccount = &dhtmlform.FormHandler{
		RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
			user := fd.GetParam("User").(*model.User)

			container := dhtml.Div().Class("border bg-light p-3").Append(
				dhtmlform.NewTextInput("displayname").Label("Display name").
					Default(user.DisplayName).Note("empty = use username: "+user.UserName),
				dhtmlform.NewPasswordInput("password1").Label("Password").Note("empty = do not change"),
				dhtmlform.NewPasswordInput("password2").Label("Confirmation"),
				dhtmlform.NewSubmitBtn().Label(mtweb.Icon("save").Label("Save")),
			)

			formBody.Append(container)
		},
		ValidateF: func(fd *dhtmlform.FormData) {
			password1 := strings.TrimSpace(fd.GetValue("password1").(string))
			password2 := strings.TrimSpace(fd.GetValue("password2").(string))

			if password1 != "" {
				if len(password1) < 6 {
					fd.SetError("password1", "Minimum password is 6 characters")
				} else {
					if password1 != password2 {
						fd.SetError("password2", "Password and confirmation do not match")
					}
				}
			}
		},
		SubmitF: func(fd *dhtmlform.FormData) {
			user := fd.GetParam("User").(*model.User)

			password := strings.TrimSpace(fd.GetValue("password1").(string))
			if password != "" {
				user.SetPassword(password)
			}

			user.DisplayName = fd.GetValue("displayname").(string)

			goapp.SaveObject(user)
		},
	}
}
