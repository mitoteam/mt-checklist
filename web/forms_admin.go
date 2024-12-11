package web

import (
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/mt-checklist/app"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mtweb"
)

func init() {
	dhtml.FormManager.Register(&dhtml.FormHandler{
		Id: "admin_checklist_edit",
		RenderF: func(form *dhtml.FormElement, fd *dhtml.FormData) {
			cl := fd.GetParam("Checklist").(*model.MtChecklist)

			form.Class("border bg-light p-3").
				Append(
					mtweb.NewFloatingFormInput("name", "text").Label("Name").DefaultValue(cl.Name),
				).
				Append(dhtml.NewFormSubmit().Label(mtweb.Icon("save").Label("Save")))
		},
		ValidateF: func(fd *dhtml.FormData) {
			name := fd.GetValue("name").(string)

			if len(name) == 0 {
				fd.SetItemError("name", "Name is required")
			}
		},
		SubmitF: func(fd *dhtml.FormData) {
			cl := fd.GetParam("Checklist").(*model.MtChecklist)

			cl.Name = fd.GetValue("name").(string)
			app.Db.Save(cl)
		},
	})
}
