package web

import (
	"slices"

	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/dhtmlform"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mttools"
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
		t := fd.GetArg("Template").(*model.Template)

		formBody.Append(dhtml.Div().Class("border bg-light p-3").Append(
			dhtmlbs.NewTextInput("name").Label("Template Name").Default(t.Name).Require(),
			dhtmlbs.NewTextInput("checklist_name").Label("Checklist Name").
				Default(t.ChecklistName).Note("Default name of created checklist"),
			dhtmlbs.NewTextarea("checklist_description").Label("Checklist Description").
				Default(t.ChecklistDescription).Note("Default description for created checklist"),
			mtweb.NewDefaultSubmitBtn(),
		))
	},
	SubmitF: func(fd *dhtmlform.FormData) {
		t := fd.GetArg("Template").(*model.Template)

		t.Name = fd.GetValue("name").(string)
		t.ChecklistName = fd.GetValue("checklist_name").(string)
		t.ChecklistDescription = fd.GetValue("checklist_description").(string)

		goapp.SaveObject(t)
	},
}

var formAdminTemplateRenumber = &dhtmlform.FormHandler{
	RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
		t := fd.GetArg("Template").(*model.Template)

		formBody.Append(dhtml.Div().Class("border bg-light p-3").Append(
			dhtml.RenderValue("Template", t.Name).Class("mb-3"),

			dhtmlbs.NewNumberInput("step").Label("Renumber step").Default(10).Require(),

			mtweb.NewDefaultSubmitBtn(),
		))
	},
	ValidateF: func(fd *dhtmlform.FormData) {
		if step, ok := mttools.AnyToInt64Ok(fd.GetValue("step")); ok {
			if step < 1 {
				fd.SetError("step", "Step should be positive number")
			}
		} else {
			fd.SetError("step", "Step should be a number")
		}
	},
	SubmitF: func(fd *dhtmlform.FormData) {
		t := fd.GetArg("Template").(*model.Template)

		step := mttools.AnyToInt64OrZero(fd.GetValue("step"))

		sortOrder := step
		for _, ti := range t.Items() {
			ti.SortOrder = sortOrder
			goapp.SaveObject(ti)

			sortOrder += step
		}

		goapp.SaveObject(t)
	},
}

var formAdminChecklistTemplateItem = &dhtmlform.FormHandler{
	RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
		item := fd.GetArg("Item").(*model.TemplateItem)

		formBody.Append(dhtml.Div().Class("border bg-light p-3").Append(
			dhtml.RenderValue("Template", item.GetTemplate().Name).Class("mb-3"),
			dhtmlbs.NewTextInput("caption").Label("Caption").Default(item.Caption).Require(),
			dhtmlbs.NewTextInput("body").Label("Body").Default(item.Body),
			NewUserSelect("responsible").Label("Responsible").Default(item.ResponsibleID),
			dhtmlbs.NewNumberInput("sort_order").Label("Sort Order").Default(item.SortOrder),
			dhtmlbs.NewNumberInput("weight").Label("Weight").Default(item.Weight),
			mtweb.NewDefaultSubmitBtn(),
		))
	},
	SubmitF: func(fd *dhtmlform.FormData) {
		item := fd.GetArg("Item").(*model.TemplateItem)

		item.Caption = fd.GetValue("caption").(string)
		item.Body = fd.GetValue("body").(string)
		item.SortOrder = mttools.AnyToInt64OrZero(fd.GetValue("sort_order"))
		item.Weight = mttools.AnyToInt64OrZero(fd.GetValue("weight"))

		if id, ok := mttools.AnyToInt64Ok(fd.GetValue("responsible")); ok {
			item.ResponsibleID = id
		}

		goapp.SaveObject(item)
	},
}

var formAdminChecklistTemplateItemDeps = &dhtmlform.FormHandler{
	RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
		item := fd.GetArg("Item").(*model.TemplateItem)
		template := item.GetTemplate()

		bodyOut := dhtml.Div().Class("border bg-light p-3").Append(
			dhtml.RenderValue("Template", template.Name).Class("mb-3"),
		)

		requiredIds := make([]int64, 0)

		for _, dep := range item.DependenciesList() {
			requiredIds = append(requiredIds, dep.RequireTemplateItemID)
		}

		for _, dItem := range template.Items() {
			if item.ID == dItem.ID {
				continue //exclude self
			}

			bodyOut.Append(
				dhtmlbs.NewCheckbox("dep_" + mttools.AnyToString(dItem.ID)).Label(dItem.Caption).
					Default(slices.Contains(requiredIds, dItem.ID)),
			)
		}

		bodyOut.Append(
			mtweb.NewDefaultSubmitBtn(),
		)

		formBody.Append(bodyOut)
	},
	SubmitF: func(fd *dhtmlform.FormData) {
		item := fd.GetArg("Item").(*model.TemplateItem)
		template := item.GetTemplate()

		goapp.Transaction(func() error {
			//delete existing ones
			goapp.DbSchema.Db().Where("template_item_id = ?", item.ID).Delete(&model.TemplateItemDependency{})

			//add checked ones
			for _, dItem := range template.Items() {
				if item.ID == dItem.ID {
					continue //exclude self
				}

				if fd.GetValue("dep_" + mttools.AnyToString(dItem.ID)).(bool) {
					dep := &model.TemplateItemDependency{}
					dep.TemplateItemID = item.ID
					dep.RequireTemplateItemID = dItem.ID

					goapp.CreateObject(dep) // TemplateItemDependency{} has no ID, so force creation rather then use SaveObject()
				}
			}

			return nil // no error = commit
		})
	},
}

var formAdminChecklist = &dhtmlform.FormHandler{
	RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
		cl := fd.GetArg("Checklist").(*model.Checklist)

		container := dhtml.Div().Class("border bg-light p-3").Append(
			dhtmlbs.NewCheckbox("active").Label("Is Active").Default(cl.IsActive),
			dhtmlbs.NewTextInput("name").Label("Name").Default(cl.Name).Require(),
			dhtmlbs.NewTextarea("description").Label("Description").Default(cl.Description),
			mtweb.NewDefaultSubmitBtn(),
		)

		formBody.Append(container)
	},
	SubmitF: func(fd *dhtmlform.FormData) {
		cl := fd.GetArg("Checklist").(*model.Checklist)

		cl.IsActive = fd.GetValue("active").(bool)
		cl.Name = fd.GetValue("name").(string)
		cl.Description = fd.GetValue("description").(string)

		goapp.SaveObject(cl)
	},
}

var formAdminChecklistItem = &dhtmlform.FormHandler{
	RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
		item := fd.GetArg("Item").(*model.ChecklistItem)

		formBody.Append(dhtml.Div().Class("border bg-light p-3").Append(
			dhtml.RenderValue("Checklist", item.GetChecklist().Name).Class("mb-3"),
			dhtmlbs.NewTextInput("caption").Label("Caption").Default(item.Caption).Require(),
			dhtmlbs.NewTextarea("body").Label("Body").Default(item.Body),
			dhtmlbs.NewNumberInput("sort_order").Label("Sort Order").Default(item.SortOrder),
			dhtmlbs.NewNumberInput("weight").Label("Weight").Default(item.Weight),
			NewUserSelect("responsible").Label("Responsible").Default(item.ResponsibleID),
			mtweb.NewDefaultSubmitBtn(),
		))
	},
	SubmitF: func(fd *dhtmlform.FormData) {
		item := fd.GetArg("Item").(*model.ChecklistItem)

		item.Caption = fd.GetValue("caption").(string)
		item.Body = fd.GetValue("body").(string)
		item.SortOrder = mttools.AnyToInt64OrZero(fd.GetValue("sort_order"))
		item.Weight = mttools.AnyToInt64OrZero(fd.GetValue("weight"))

		if id, ok := mttools.AnyToInt64Ok(fd.GetValue("responsible")); ok {
			item.ResponsibleID = id
		}

		goapp.SaveObject(item)
	},
}

var formAdminChecklistItemDeps = &dhtmlform.FormHandler{
	RenderF: func(formBody *dhtml.HtmlPiece, fd *dhtmlform.FormData) {
		ci := fd.GetArg("Item").(*model.ChecklistItem)
		cl := ci.GetChecklist()

		bodyOut := dhtml.Div().Class("border bg-light p-3").Append(
			dhtml.RenderValue("Checklist", cl.Name).Class("mb-3"),
		)

		requiredIds := make([]int64, 0)

		for _, dep := range ci.DependenciesList() {
			requiredIds = append(requiredIds, dep.RequireChecklistItemID)
		}

		for _, dItem := range cl.Items() {
			if ci.ID == dItem.ID {
				continue //exclude self
			}

			bodyOut.Append(
				dhtmlbs.NewCheckbox("dep_" + mttools.AnyToString(dItem.ID)).Label(dItem.Caption).
					Default(slices.Contains(requiredIds, dItem.ID)),
			)
		}

		bodyOut.Append(
			mtweb.NewDefaultSubmitBtn(),
		)

		formBody.Append(bodyOut)
	},
	SubmitF: func(fd *dhtmlform.FormData) {
		ci := fd.GetArg("Item").(*model.ChecklistItem)
		cl := ci.GetChecklist()

		goapp.Transaction(func() error {
			//delete existing ones
			goapp.DbSchema.Db().Where("checklist_item_id = ?", ci.ID).Delete(&model.ChecklistItemDependency{})

			//add checked ones
			for _, dItem := range cl.Items() {
				if ci.ID == dItem.ID {
					continue //exclude self
				}

				if fd.GetValue("dep_" + mttools.AnyToString(dItem.ID)).(bool) {
					dep := &model.ChecklistItemDependency{}
					dep.ChecklistItemID = ci.ID
					dep.RequireChecklistItemID = dItem.ID

					goapp.CreateObject(dep)
				}
			}

			return nil // no error = commit
		})
	},
}
