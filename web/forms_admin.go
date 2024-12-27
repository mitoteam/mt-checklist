package web

import (
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlform"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mtweb"
)

var formAdminUserEdit = &dhtmlform.FormHandler{
	RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
		user := fd.GetArg("User").(*model.User)

		container := dhtml.Div().Class("border bg-light p-3").Append(
			dhtmlform.NewTextInput("username").Label("Username").Require().Default(user.UserName),
			dhtmlform.NewTextInput("displayname").Label("Display name").Default(user.DisplayName),
			dhtmlform.NewCheckbox("is_active").Label("Active").Default(user.IsActive),
			dhtmlform.NewSubmitBtn().Label(mtweb.Icon("save").Label("Save")),
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
			dhtmlform.NewPasswordInput("password1").Label("Password"),
			dhtmlform.NewPasswordInput("password2").Label("Confirmation"),
			dhtmlform.NewSubmitBtn().Label(mtweb.Icon("save").Label("Save")),
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

func init() {
	Forms.AdminChecklist = &dhtmlform.FormHandler{
		RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
			cl := fd.GetParam("Checklist").(*model.Checklist)

			container := dhtml.Div().Class("border bg-light p-3").Append(
				dhtmlform.NewTextInput("name").Label("Name").Default(cl.Name).Require(),
				dhtmlform.NewSubmitBtn().Label(mtweb.Icon("save").Label("Save")),
			)

			formBody.Append(container)
		},
		SubmitF: func(fd *dhtmlform.FormData) {
			cl := fd.GetParam("Checklist").(*model.Checklist)

			cl.Name = fd.GetValue("name").(string)

			goapp.SaveObject(cl)
		},
	}

	Forms.AdminChecklistTemplate = &dhtmlform.FormHandler{
		RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
			t := fd.GetParam("Template").(*model.ChecklistTemplate)

			formBody.Append(dhtml.Div().Class("border bg-light p-3").Append(
				dhtmlform.NewTextInput("name").Label("Template Name").Default(t.Name).Require(),
				dhtmlform.NewTextInput("checklist_name").Label("Checklist Name").
					Default(t.ChecklistName).Note("Default name of created checklist"),
				dhtmlform.NewSubmitBtn().Label(mtweb.Icon("save").Label("Save")),
			))
		},
		SubmitF: func(fd *dhtmlform.FormData) {
			t := fd.GetParam("Template").(*model.ChecklistTemplate)

			t.Name = fd.GetValue("name").(string)
			t.ChecklistName = fd.GetValue("checklist_name").(string)

			goapp.SaveObject(t)
		},
	}

	Forms.AdminChecklistTemplateItem = &dhtmlform.FormHandler{
		RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
			item := fd.GetParam("Item").(*model.ChecklistTemplateItem)

			formBody.Append(dhtml.Div().Class("border bg-light p-3").Append(
				dhtml.RenderValue("Template", item.GetChecklistTemplate().Name).Class("mb-3"),
				dhtmlform.NewTextInput("caption").Label("Caption").Default(item.Caption).Require(),
				dhtmlform.NewTextInput("body").Label("Body").Default(item.Body),
				dhtmlform.NewSubmitBtn().Label(mtweb.Icon("save").Label("Save")),
			))
		},
		SubmitF: func(fd *dhtmlform.FormData) {
			item := fd.GetParam("Item").(*model.ChecklistTemplateItem)

			item.Caption = fd.GetValue("caption").(string)
			item.Body = fd.GetValue("body").(string)

			item.ResponsibleID = fd.GetArg("User").(*model.User).ID

			goapp.SaveObject(item)
		},
	}
}
