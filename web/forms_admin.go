package web

import (
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/goappbase"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mtweb"
)

func init() {
	dhtml.FormManager.Register(&dhtml.FormHandler{
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
	})

	dhtml.FormManager.Register(&dhtml.FormHandler{
		Id: "admin_user_edit",
		RenderF: func(form *dhtml.FormElement, fd *dhtml.FormData) {
			user := fd.GetParam("User").(*model.User)

			form.Class("border bg-light p-3").Append(
				dhtml.NewFormInput("username", "text").Label("Username").DefaultValue(user.UserName),
			).Append(
				dhtml.NewFormInput("displayname", "text").Label("Display name").DefaultValue(user.DisplayName),
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

			goappbase.SaveObject(user)
		},
	})
}
