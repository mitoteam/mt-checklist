package web

import (
	"github.com/gin-contrib/sessions"
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/mt-checklist/app"
	"github.com/mitoteam/mtweb"
)

func init() {
	dhtml.FormManager.Register(&dhtml.FormHandler{
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
				user := app.AuthorizeUser(username, password)

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
	})
}
