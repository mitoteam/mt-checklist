package web

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/dhtmlform"
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-checklist/app"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mtweb"
)

type PageBuilder struct {
	ctx *mbr.MbrContext

	// "title" = H1 page title
	// "main" = main content
	regions dhtml.NamedHtmlPieces
}

func PageBuilderRouteHandler(buildPageF func(*PageBuilder) any) func(ctx *mbr.MbrContext) any {
	return func(ctx *mbr.MbrContext) any {
		p := &PageBuilder{
			regions: dhtml.NewNamedHtmlPieces(),
			ctx:     ctx,
		}

		out := buildPageF(p)

		if err, ok := out.(error); ok {
			return err
		}

		if p.HasMain() {
			ctx.Request().Header.Add("Content-Type", "text/html;charset=utf-8")
			return p.String()
		} else {
			return out
		}
	}
}

func (p *PageBuilder) User() (user *model.User) {
	if v, ok := p.ctx.GetOk("User"); ok {
		user = v.(*model.User)
	}

	return user
}

func (p *PageBuilder) Title(v any) *PageBuilder {
	p.regions.Add("title", v)
	return p
}

func (p *PageBuilder) GetTitle() *dhtml.HtmlPiece {
	return p.regions.Get("title")
}

func (p *PageBuilder) Main(v any) *PageBuilder {
	p.regions.Add("main", v)
	return p
}

func (p *PageBuilder) GetMain() *dhtml.HtmlPiece {
	return p.regions.Get("main")
}

func (p *PageBuilder) HasMain() bool {
	return !p.regions.IsEmpty("main")
}

func (p *PageBuilder) String() string {
	return p.render().String()
}

func (p *PageBuilder) render() (out *dhtml.HtmlPiece) {
	document := dhtml.NewHtmlDocument()

	var head_title = app.App.AppName

	title := p.regions.Get("title")
	if !title.IsEmpty() {
		head_title = title.String() + " | " + head_title
	}

	document.
		Charset("utf-8").
		Title(head_title).
		Icon("/favicon.ico").
		Stylesheet("/assets/vendor/bootstrap.min.css").
		Stylesheet("/assets/vendor/fontawesome.min.css").
		Stylesheet("/assets/vendor/regular.min.css").
		Stylesheet("/assets/css/style.css")

	container := dhtml.Div().Class("container my-3")

	container.Append(p.renderHeader())

	// H1 page title
	if !title.IsEmpty() {
		container.Append(dhtml.NewTag("h1").Append(title))
	}

	container.Append(dhtml.Div().Class("region-main").Append(p.GetMain()))

	container.Append(p.renderFooter())

	document.Body().Append(container)

	//scripts
	document.Body().
		Append(dhtml.NewTag("script").Attribute("src", "/assets/vendor/bootstrap.bundle.min.js")).
		//time for vue has not come yet
		//Append(dhtml.NewTag("script").Attribute("src", "/assets/vendor/vue.global.prod.js")).
		Append(dhtml.NewTag("script").Attribute("src", "/assets/script.min.js"))

	return dhtml.Piece(document)
}

func (p *PageBuilder) renderHeader() (out dhtml.HtmlPiece) {
	user := p.User()

	header := dhtml.Div().Class("region-header border bg-light p-3 mb-3").Attribute("role", "header")

	header_left := dhtml.Div().
		Append(dhtml.Div().Append(dhtml.NewLink(mbr.Url(RootCtl.Home)).Label(app.App.AppName).Class("text-decoration-none"))).
		Append(dhtml.Div().Class("small text-muted").Append("v." + app.App.Version))

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
				Append(dhtml.NewLink(mbr.Url(RootCtl.MyAccount, "destination", "/")).Label(mtweb.Icon("cog"))),
		)
	}

	header.Append(dhtmlbs.NewJustifiedLR().L(header_left).R(header_right))

	out.Append(header)
	return out
}

func (p *PageBuilder) renderFooter() (out dhtml.HtmlPiece) {
	out.Append(dhtml.Div().Class("region-footer border bg-light p-3 mt-3").Append(
		dhtmlbs.NewJustifiedLR().
			L(fmt.Sprintf("%s v.%s", app.App.AppName, app.App.Version)).
			R(
				dhtml.Div().Class("small text-muted").
					Text("by ").
					Append(dhtml.NewLink("https://www.mito-team.com").Label("MiTo Team").Target("blank")),
			),
	))
	return out
}

// Builds new dhtml.FormContext to be used with form builder
func (p *PageBuilder) FormContext() *dhtmlform.FormContext {
	fc := dhtmlform.NewFormContext(p.ctx.Writer(), p.ctx.Request())

	// some useful for every form things
	fc.SetParam("MbrContext", p.ctx)
	fc.SetParam("User", p.User())

	//default redirect from "destination" query parameter
	if destination := p.ctx.Request().URL.Query().Get("destination"); destination != "" {
		fc.SetRedirect(destination)
	}

	return fc
}

// ===================== PageBuilderOLD =================================

type PageBuilderOLD struct {
	context *gin.Context

	// "title" = H1 page title
	// "main" = main content
	regions dhtml.NamedHtmlPieces
}

func NewPageBuilderOLD(context *gin.Context) *PageBuilderOLD {
	p := &PageBuilderOLD{
		regions: dhtml.NewNamedHtmlPieces(),
		context: context,
	}

	return p
}

func (p *PageBuilderOLD) GetGinContext() *gin.Context {
	return p.context
}

func (p *PageBuilderOLD) GetSession() sessions.Session {
	return sessions.Default(p.context)
}

// Builds new dhtml.FormContext to be used with form builder
func (p *PageBuilderOLD) FormContext() *dhtmlform.FormContext {
	fc := dhtmlform.NewFormContext(p.context.Writer, p.context.Request)

	//current user from session if he authorized
	if user, ok := p.context.Get("User"); ok {
		fc.SetArg("User", user.(*model.User))
	}

	return fc
}

func (p *PageBuilderOLD) Title(v any) *PageBuilderOLD {
	p.regions.Add("title", v)
	return p
}

func (p *PageBuilderOLD) GetTitle() *dhtml.HtmlPiece {
	return p.regions.Get("title")
}

func (p *PageBuilderOLD) Main(v any) *PageBuilderOLD {
	p.regions.Add("main", v)
	return p
}

func (p *PageBuilderOLD) GetMain() *dhtml.HtmlPiece {
	return p.regions.Get("main")
}

func (p *PageBuilderOLD) String() string {
	return p.render().String()
}

func (p *PageBuilderOLD) User() (user *model.User) {
	if v, ok := p.context.Get("User"); ok {
		user = v.(*model.User)
	}

	return user
}

func (p *PageBuilderOLD) render() (out *dhtml.HtmlPiece) {
	document := dhtml.NewHtmlDocument()

	var head_title = app.App.AppName

	title := p.regions.Get("title")
	if !title.IsEmpty() {
		head_title = title.String() + " | " + head_title
	}

	document.
		Charset("utf-8").
		Title(head_title).
		Icon("/favicon.ico").
		Stylesheet("/assets/vendor/bootstrap.min.css").
		Stylesheet("/assets/vendor/fontawesome.min.css").
		Stylesheet("/assets/vendor/regular.min.css").
		Stylesheet("/assets/css/style.css")

	container := dhtml.Div().Class("container my-3")

	container.Append(p.renderHeader())

	// H1 page title
	if !title.IsEmpty() {
		container.Append(dhtml.NewTag("h1").Append(title))
	}

	container.Append(dhtml.Div().Class("region-main").Append(p.GetMain()))

	container.Append(p.renderFooter())

	document.Body().Append(container)

	//scripts
	document.Body().
		Append(dhtml.NewTag("script").Attribute("src", "/assets/vendor/bootstrap.bundle.min.js")).
		//time for vue has not come yet
		//Append(dhtml.NewTag("script").Attribute("src", "/assets/vendor/vue.global.prod.js")).
		Append(dhtml.NewTag("script").Attribute("src", "/assets/script.min.js"))

	return dhtml.Piece(document)
}

func (p *PageBuilderOLD) renderHeader() (out dhtml.HtmlPiece) {
	user := p.User()

	header := dhtml.Div().Class("region-header border bg-light p-3 mb-3").Attribute("role", "header")

	header_left := dhtml.Div().
		Append(dhtml.Div().Append(dhtml.NewLink("/").Label(app.App.AppName).Class("text-decoration-none"))).
		Append(dhtml.Div().Class("small text-muted").Append("v." + app.App.Version))

	header_right := dhtml.Div().Class("text-end")

	if user != nil {
		header_right.Append(dhtml.Div().
			Text(user.GetDisplayName()).
			Append(
				dhtml.NewLink("/logout").Label(mtweb.Icon("arrow-right-from-bracket")).
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
				Append(dhtml.NewLink("/account").Label(mtweb.Icon("cog"))),
		)
	}

	header.Append(dhtmlbs.NewJustifiedLR().L(header_left).R(header_right))

	out.Append(header)
	return out
}

func (p *PageBuilderOLD) renderFooter() (out dhtml.HtmlPiece) {
	out.Append(dhtml.Div().Class("region-footer border bg-light p-3 mt-3").Append(
		dhtmlbs.NewJustifiedLR().
			L(fmt.Sprintf("%s v.%s", app.App.AppName, app.App.Version)).
			R(
				dhtml.Div().Class("small text-muted").
					Text("by ").
					Append(dhtml.NewLink("https://www.mito-team.com").Label("MiTo Team").Target("blank")),
			),
	))
	return out
}
