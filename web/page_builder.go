package web

import (
	"fmt"

	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/dhtmlform"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-checklist/app"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mtweb"
)

type PageBuilder struct {
	mtweb.PageBuilderBase
}

func PageBuilderRouteHandler(buildPageF func(*PageBuilder) any) func(ctx *mbr.MbrContext) any {
	return func(ctx *mbr.MbrContext) any {
		// set up page builder
		p := &PageBuilder{
			PageBuilderBase: mtweb.NewPageBuilderBase(ctx),
		}

		p.BuildHeadTitleF = func() string {
			title := app.Options.SiteName()

			if p.GetTitle() != "" {
				title = p.GetTitle() + " | " + title
			}

			return title
		}

		p.RenderF = p.render

		// render page content
		out := buildPageF(p)

		if err, ok := out.(error); ok { // error happen, return it as-is
			return err
		} else if p.Ctx.IsRedirect() { //redirect is already set, so we don't need to return content
			return nil
		} else { // render the page as html
			html, err := p.Render()

			if err != nil {
				return err
			}

			return html
		}
	}
}

func (p *PageBuilder) User() (user *model.User) {
	if v, ok := p.Ctx.GetOk("User"); ok {
		user = v.(*model.User)
	}

	return user
}

// "override" NewFormContext to make additional data available to form builder
func (p *PageBuilder) NewFormContext() *dhtmlform.FormContext {
	fc := p.PageBuilderBase.NewFormContext()

	// make current use available to form builder
	fc.SetParam("User", p.User())

	return fc
}

func (p *PageBuilder) render() error {
	document := p.GetDocument().
		Icon("/favicon.ico").
		Stylesheet("/assets/css/style.css")

	container := dhtml.Div().Class("container my-3")

	container.Append(p.renderHeader())

	// H1 page title
	if p.GetTitle() != "" {
		container.Append(dhtml.NewTag("h1").Append(p.GetTitle()))
	}

	container.Append(dhtml.Div().Class("region-main").Append(p.GetMain()))

	container.Append(p.renderFooter())

	document.Body().Append(container)

	//scripts
	document.Body().
		Append(dhtml.NewTag("script").Attribute("src", "/assets/script.min.js"))

	return nil
}

func (p *PageBuilder) renderHeader() (out dhtml.HtmlPiece) {
	user := p.User()

	header := dhtml.Div().Class("region-header border bg-light p-3 mb-3").Attribute("role", "header")

	header_left := dhtml.Div().
		Append(dhtml.Div().Append(dhtml.NewLink(mbr.Url(RootCtl.Home)).Label(app.Options.SiteName()).Class("text-decoration-none")))

	motto := app.Options.SiteMotto()
	if motto != "" {
		header_left.Append(dhtml.Div().Class("small text-muted").Append(motto))
	}

	header_right := dhtml.Div().Class("text-end")

	if user != nil {
		header_right.Append(dhtml.Div().
			Text(user.GetDisplayName()).
			Append(
				dhtml.NewLink(mbr.Url(RootCtl.Logout)).Label(mtweb.Icon("arrow-right-from-bracket")).
					Class("ms-1").Title("Sign Out"),
			))

		var icon *mtweb.IconElement

		if user.IsAdmin() {
			icon = mtweb.Icon("user-police-tie").Title("administrator")
		} else {
			icon = mtweb.Icon("user").Title("user")
		}

		icon.Label(user.UserName)

		header_right.Append(
			dhtml.Div().Class("text-muted").
				Append(icon).
				Append(dhtml.NewLink(mbr.Url(RootCtl.MyAccount, "destination", "/")).Label(mtweb.Icon(mtweb.FaIconOptions))),
		)
	}

	header.Append(dhtmlbs.NewJustifiedLR().L(header_left).R(header_right))

	out.Append(header)
	return out
}

func (p *PageBuilder) renderFooter() (out dhtml.HtmlPiece) {
	out.Append(dhtml.Div().Class("region-footer border bg-light p-3 mt-3").Append(
		dhtmlbs.NewJustifiedLR().
			L(
				fmt.Sprintf("This instance: v%s", app.App.Version),
				dhtml.Span().Class("small text-muted ms-2").Append(
					mtweb.Icon(mtweb.IconTimestamp).Label(app.App.BuildTime),
				),
				dhtml.Div().Append(
					dhtml.NewLink("https://github.com/mitoteam/mt-checklist").Label(
						dhtml.UnsafeText("<img alt=\"GitHub Release\" src=\"https://img.shields.io/github/v/release/mitoteam/mt-checklist?style=flat-square&logo=github&label=latest%20version\">"),
					),
				),
				dhtml.Div().Class("mt-1").Append(
					mtweb.NewSmBtn(
						"https://github.com/mitoteam/mt-checklist/issues/new?template=bug_report.md", "bug",
					).Label(dhtml.Span().Class("ms-1").Append("Report a Bug")).Target("blank"),
					mtweb.NewSmBtn(
						"https://github.com/mitoteam/mt-checklist/issues/new?template=feature_request.md", "lightbulb-on",
					).Label(dhtml.Span().Class("ms-1").Append("Suggest a Feature")).Target("blank"),
				),
			).
			R(
				dhtml.Div().Class("small text-end").Append(
					app.App.AppName+" by ",
					dhtml.NewLink("https://www.mito-team.com").Label("MiTo Team").Target("blank"),
				),
				dhtml.Div().Class("small text-muted text-end").Append(goapp.MOTTO),
			),
	))
	return out
}
