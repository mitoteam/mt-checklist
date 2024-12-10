package web

import (
	"github.com/mitoteam/dhtml"
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
			if v, ok := fd.GetValue("weha").(string); ok {
				if len(v) < 3 {
					fd.SetItemError("weha", "At least three characters expected")
				}
			}
		},
		SubmitF: func(fd *dhtml.FormData) {
			fd.SetRedirect("/")
		},
	})
}
