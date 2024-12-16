package web

import (
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/goappbase"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mtweb"
)

func init() {
	Forms.AdminChecklist = &dhtml.FormHandler{
		Id: "admin_checklist_edit",
		RenderF: func(form *dhtml.FormElement, fd *dhtml.FormData) {
			cl := fd.GetParam("Checklist").(*model.Checklist)

			form.Class("border bg-light p-3").
				Append(
					mtweb.NewFloatingFormInput("name", "text").Label("Name").DefaultValue(cl.Name),
				).
				Append(dhtml.NewFormSubmit().Label(mtweb.Icon("save").Label("Save")))
		},
		ValidateF: func(fd *dhtml.FormData) {
			if len(fd.GetValue("name").(string)) == 0 {
				fd.SetItemError("name", "Name is required")
			}
		},
		SubmitF: func(fd *dhtml.FormData) {
			cl := fd.GetParam("Checklist").(*model.Checklist)

			cl.Name = fd.GetValue("name").(string)

			goappbase.SaveObject(cl)
		},
	}
	dhtml.FormManager.Register(Forms.AdminChecklist)

	Forms.AdminUserEdit = &dhtml.FormHandler{
		Id: "admin_user_edit",
		RenderF: func(form *dhtml.FormElement, fd *dhtml.FormData) {
			user := fd.GetParam("User").(*model.User)

			form.Class("border bg-light p-3").Append(
				dhtml.NewFormInput("username", "text").Label("Username").DefaultValue(user.UserName),
			).Append(
				dhtml.NewFormInput("displayname", "text").Label("Display name").DefaultValue(user.DisplayName),
			).Append(
				dhtml.NewFormCheckbox("is_active").Label("Active").DefaultValue(user.IsActive),
			).Append(
				dhtml.NewFormSubmit().Label(mtweb.Icon("save").Label("Save")),
			)
		},
		ValidateF: func(fd *dhtml.FormData) {
			if len(fd.GetValue("username").(string)) == 0 {
				fd.SetItemError("username", "Username is required")
			}

			if len(fd.GetValue("displayname").(string)) == 0 {
				fd.SetValue("displayname", fd.GetValue("username"))
			}
		},
		SubmitF: func(fd *dhtml.FormData) {
			user := fd.GetParam("User").(*model.User)

			user.UserName = fd.GetValue("username").(string)
			user.DisplayName = fd.GetValue("displayname").(string)
			user.IsActive = fd.GetValue("is_active").(bool)

			goappbase.SaveObject(user)
		},
	}
	dhtml.FormManager.Register(Forms.AdminUserEdit)

	Forms.AdminUserPassword = &dhtml.FormHandler{
		Id: "admin_user_password",
		RenderF: func(form *dhtml.FormElement, fd *dhtml.FormData) {
			form.Class("border bg-light p-3").Append(
				dhtml.NewFormInput("password1", "password").Label("Password"),
			).Append(
				dhtml.NewFormInput("password2", "password").Label("Confirmation"),
			).Append(
				dhtml.NewFormSubmit().Label(mtweb.Icon("save").Label("Save")),
			)
		},
		ValidateF: func(fd *dhtml.FormData) {
			password1 := fd.GetValue("password1").(string)
			password2 := fd.GetValue("password2").(string)

			if len(password1) < 6 {
				fd.SetItemError("password1", "Minimum password is 6 characters")
			} else {
				if password1 != password2 {
					fd.SetItemError("password2", "Password and confirmation do not match")
				}
			}
		},
		SubmitF: func(fd *dhtml.FormData) {
			user := fd.GetParam("User").(*model.User)

			user.SetPassword(fd.GetValue("password1").(string))

			goappbase.SaveObject(user)
		},
	}
	dhtml.FormManager.Register(Forms.AdminUserPassword)

	Forms.AdminChecklistTemplate = &dhtml.FormHandler{
		Id: "admin_checklist_template",
		RenderF: func(form *dhtml.FormElement, fd *dhtml.FormData) {
			t := fd.GetParam("Template").(*model.ChecklistTemplate)

			form.Class("border bg-light p-3").Append(
				dhtml.NewFormInput("name", "text").Label("Template Name").DefaultValue(t.Name),
			).Append(
				dhtml.NewFormInput("checklist_name", "text").Label("Checklist Name").
					DefaultValue(t.ChecklistName).Note("Default name of created checklist"),
			).Append(
				dhtml.NewFormSubmit().Label(mtweb.Icon("save").Label("Save")),
			)
		},
		ValidateF: func(fd *dhtml.FormData) {
			name := fd.GetValue("name").(string)

			if len(name) == 0 {
				fd.SetItemError("name", "Template Name is required")
			}
		},
		SubmitF: func(fd *dhtml.FormData) {
			t := fd.GetParam("Template").(*model.ChecklistTemplate)

			t.Name = fd.GetValue("name").(string)
			t.ChecklistName = fd.GetValue("checklist_name").(string)

			goappbase.SaveObject(t)
		},
	}
	dhtml.FormManager.Register(Forms.AdminChecklistTemplate)

	Forms.AdminChecklistTemplateItem = &dhtml.FormHandler{
		Id: "admin_checklist_template_item",
		RenderF: func(form *dhtml.FormElement, fd *dhtml.FormData) {
			item := fd.GetParam("Item").(*model.ChecklistTemplateItem)

			form.Append(dhtml.RenderValue("Template", item.GetChecklistTemplate().Name).Class("mb-3"))

			form.Class("border bg-light p-3").Append(
				dhtml.NewFormInput("caption", "text").Label("Caption").DefaultValue(item.Caption),
			).Append(
				dhtml.NewFormSubmit().Label(mtweb.Icon("save").Label("Save")),
			)
		},
		ValidateF: func(fd *dhtml.FormData) {
			caption := fd.GetValue("caption").(string)

			if len(caption) == 0 {
				fd.SetItemError("name", "Caption is required")
			}
		},
		SubmitF: func(fd *dhtml.FormData) {
			item := fd.GetParam("Item").(*model.ChecklistTemplateItem)

			item.Caption = fd.GetValue("caption").(string)

			goappbase.SaveObject(item)
		},
	}
	dhtml.FormManager.Register(Forms.AdminChecklistTemplateItem)
}
