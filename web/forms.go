package web

import (
	"strings"

	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/dhtmlform"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mtweb"
)

var formLogin = &dhtmlform.FormHandler{
	RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
		formBody.Append(dhtml.Div().Class("border bg-light p-3").Append(
			dhtmlbs.NewFloatingTextInput("username").Label("Username").Require(),
			dhtmlbs.NewFloatingPasswordInput("password").Label("Password").Require(),

			mtweb.NewIconSubmitBtn("arrow-right-to-bracket", "Sign In"),
		))
	},
	ValidateF: func(fd *dhtmlform.FormData) {
		if !fd.HasError() {
			username := fd.GetValue("username").(string)
			password := fd.GetValue("password").(string)

			user := model.AuthorizeUser(username, password)

			ctx := fd.GetParam("MbrContext").(*mbr.MbrContext)

			if user == nil {
				session := Session(ctx.Request())
				delete(session.Values, sessionIdField) //remove old value if it was set
				session.Save(ctx.Request(), ctx.Writer())

				fd.SetError("", "User not found or wrong password given")
			} else {
				fd.SetParam(sessionIdField, user.SessionId)
			}
		}
	},
	SubmitF: func(fd *dhtmlform.FormData) {
		ctx := fd.GetParam("MbrContext").(*mbr.MbrContext)
		session := Session(ctx.Request())

		session.Values[sessionIdField] = fd.GetParam(sessionIdField).(string)
		session.Save(ctx.Request(), ctx.Writer())
	},
}

var formMyAccount = &dhtmlform.FormHandler{
	RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
		user := fd.GetParam("User").(*model.User)

		container := dhtml.Div().Class("border bg-light p-3").Append(
			dhtmlbs.NewTextInput("displayname").Label("Display name").
				Default(user.DisplayName).Note("empty = use username: "+user.UserName),
			dhtmlbs.NewPasswordInput("password1").Label("Password").Note("empty = do not change"),
			dhtmlbs.NewPasswordInput("password2").Label("Password confirmation"),
			mtweb.NewDefaultSubmitBtn(),
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
