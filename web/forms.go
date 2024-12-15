package web

import (
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/goappbase"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mtweb"
)

// helper type to have all form handlers in one place
type formsType struct {
	Login *dhtml.FormHandler

	AdminUserEdit     *dhtml.FormHandler
	AdminUserPassword *dhtml.FormHandler

	AdminChecklist *dhtml.FormHandler

	MyAccount *dhtml.FormHandler
}

var Forms formsType

func init() {
	Forms.Login = &dhtml.FormHandler{
		Id: "login",
		RenderF: func(form *dhtml.FormElement, fd *dhtml.FormData) {
			form.Class("border bg-light p-3").
				Append(
					mtweb.NewFloatingFormInput("username", "text").Placeholder("Username").Label("Username"),
				).
				Append(
					mtweb.NewFloatingFormInput("password", "password").Label("Password"),
				).
				Append(dhtml.NewFormSubmit().Label(mtweb.Icon("arrow-right-to-bracket").Label("Sign In")))
		},
		ValidateF: func(fd *dhtml.FormData) {
			username := fd.GetValue("username").(string)
			password := fd.GetValue("password").(string)

			if len(username) == 0 {
				fd.SetItemError("username", "Username required")
			}

			if len(password) == 0 {
				fd.SetItemError("password", "Password required")
			}

			if !fd.HasError() {
				user := model.AuthorizeUser(username, password)

				if user != nil {
					fd.SetValue("userID", user.ID)
				} else {
					if session, ok := fd.GetParam("Session").(sessions.Session); ok {
						session.Delete("userID")
						session.Save()
					}

					fd.SetError("User not found or wrong password given")
				}
			}
		},
		SubmitF: func(fd *dhtml.FormData) {
			if session, ok := fd.GetParam("Session").(sessions.Session); ok {

				session.Set("userID", fd.GetValue("userID").(int64))
				session.Save()
			}
		},
	}
	dhtml.FormManager.Register(Forms.Login)

	Forms.MyAccount = &dhtml.FormHandler{
		Id: "my_account",
		RenderF: func(form *dhtml.FormElement, fd *dhtml.FormData) {
			user := fd.GetParam("User").(*model.User)

			form.Class("border bg-light p-3").Append(
				dhtml.NewFormInput("displayname", "text").Label("Display name").
					DefaultValue(user.DisplayName).Note("empty = use username: " + user.UserName),
			).Append(
				dhtml.NewFormInput("password1", "password").Label("Password").Note("empty = do not change"),
			).Append(
				dhtml.NewFormInput("password2", "password").Label("Confirmation"),
			).Append(
				dhtml.NewFormSubmit().Label(mtweb.Icon("save").Label("Save")),
			)
		},
		ValidateF: func(fd *dhtml.FormData) {
			password1 := strings.TrimSpace(fd.GetValue("password1").(string))
			password2 := strings.TrimSpace(fd.GetValue("password2").(string))

			if password1 != "" {
				if len(password1) < 6 {
					fd.SetItemError("password1", "Minimum password is 6 characters")
				} else {
					if password1 != password2 {
						fd.SetItemError("password2", "Password and confirmation do not match")
					}
				}
			}
		},
		SubmitF: func(fd *dhtml.FormData) {
			user := fd.GetParam("User").(*model.User)

			password := strings.TrimSpace(fd.GetValue("password1").(string))
			if password != "" {
				user.SetPassword(password)
			}

			user.DisplayName = fd.GetValue("displayname").(string)

			goappbase.SaveObject(user)
		},
	}
	dhtml.FormManager.Register(Forms.MyAccount)
}
