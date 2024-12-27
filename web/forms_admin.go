package web

import (
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/dhtmlform"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mtweb"
)

var formAdminUserEdit = &dhtmlform.FormHandler{
	RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
		user := fd.GetArg("User").(*model.User)

		container := dhtml.Div().Class("border bg-light p-3").Append(
			dhtmlbs.NewTextInput("username").Label("Username").Require().Default(user.UserName),
			dhtmlbs.NewTextInput("displayname").Label("Display name").Default(user.DisplayName),
			dhtmlbs.NewCheckbox("is_active").Label("Active").Default(user.IsActive).Note("Uncheck to disable sign-in"),
			mtweb.NewDefaultSubmitBtn(),
		)

		formBody.Append(container)
	},
	ValidateF: func(fd *dhtmlform.FormData) {
		// get display name from username if not set
		if len(fd.GetValue("displayname").(string)) == 0 {
			fd.SetControlValue("displayname", fd.GetValue("username"))
		}
	},
	SubmitF: func(fd *dhtmlform.FormData) {
		user := fd.GetArg("User").(*model.User)

		user.UserName = fd.GetValue("username").(string)
		user.DisplayName = fd.GetValue("displayname").(string)
		user.IsActive = fd.GetValue("is_active").(bool)

		goapp.SaveObject(user)
	},
}

var formAdminUserPassword = &dhtmlform.FormHandler{
	RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
		container := dhtml.Div().Class("border bg-light p-3").Append(
			dhtmlbs.NewPasswordInput("password1").Label("Password"),
			dhtmlbs.NewPasswordInput("password2").Label("Confirmation"),
			mtweb.NewDefaultSubmitBtn(),
		)

		formBody.Append(container)
	},
	ValidateF: func(fd *dhtmlform.FormData) {
		password1 := fd.GetValue("password1").(string)
		password2 := fd.GetValue("password2").(string)

		if len(password1) < 6 {
			fd.SetError("password1", "Minimum password is 6 characters")
		} else {
			if password1 != password2 {
				fd.SetError("password2", "Password and confirmation do not match")
			}
		}
	},
	SubmitF: func(fd *dhtmlform.FormData) {
		user := fd.GetArg("User").(*model.User)

		user.SetPassword(fd.GetValue("password1").(string))

		goapp.SaveObject(user)
	},
}

var formAdminChecklistTemplate = &dhtmlform.FormHandler{
	RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
		t := fd.GetArg("Template").(*model.ChecklistTemplate)

		formBody.Append(dhtml.Div().Class("border bg-light p-3").Append(
			dhtmlbs.NewTextInput("name").Label("Template Name").Default(t.Name).Require(),
			dhtmlbs.NewTextInput("checklist_name").Label("Checklist Name").
				Default(t.ChecklistName).Note("Default name of created checklist"),
			mtweb.NewDefaultSubmitBtn(),
		))
	},
	SubmitF: func(fd *dhtmlform.FormData) {
		t := fd.GetArg("Template").(*model.ChecklistTemplate)

		t.Name = fd.GetValue("name").(string)
		t.ChecklistName = fd.GetValue("checklist_name").(string)

		goapp.SaveObject(t)
	},
}

var formAdminChecklistTemplateItem = &dhtmlform.FormHandler{
	RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
		item := fd.GetArg("Item").(*model.ChecklistTemplateItem)

		formBody.Append(dhtml.Div().Class("border bg-light p-3").Append(
			dhtml.RenderValue("Template", item.GetChecklistTemplate().Name).Class("mb-3"),
			dhtmlbs.NewTextInput("caption").Label("Caption").Default(item.Caption).Require(),
			dhtmlbs.NewTextInput("body").Label("Body").Default(item.Body),
			mtweb.NewDefaultSubmitBtn(),
		))
	},
	SubmitF: func(fd *dhtmlform.FormData) {
		item := fd.GetArg("Item").(*model.ChecklistTemplateItem)

		item.Caption = fd.GetValue("caption").(string)
		item.Body = fd.GetValue("body").(string)

		//current user
		item.ResponsibleID = fd.GetParam("User").(*model.User).ID

		goapp.SaveObject(item)
	},
}

var formAdminChecklist = &dhtmlform.FormHandler{
	RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
		cl := fd.GetArg("Checklist").(*model.Checklist)

		container := dhtml.Div().Class("border bg-light p-3").Append(
			dhtmlbs.NewTextInput("name").Label("Name").Default(cl.Name).Require(),
			mtweb.NewDefaultSubmitBtn(),
		)

		formBody.Append(container)
	},
	SubmitF: func(fd *dhtmlform.FormData) {
		cl := fd.GetArg("Checklist").(*model.Checklist)

		cl.Name = fd.GetValue("name").(string)

		goapp.SaveObject(cl)
	},
}
