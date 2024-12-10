package web

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/mt-checklist/app"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mtweb"
)

type PageTemplate struct {
	context *gin.Context

	// "title" = H1 page title
	// "main" = main content
	regions dhtml.NamedHtmlPieces
}

func NewPageTemplate(context *gin.Context) *PageTemplate {
	p := &PageTemplate{
		regions: dhtml.NewNamedHtmlPieces(),
		context: context,
	}

	return p
}

func (p *PageTemplate) GetContext() *gin.Context {
	return p.context
}

func (p *PageTemplate) Title(v any) *PageTemplate {
	p.regions.Add("title", v)
	return p
}

func (p *PageTemplate) GetTitle() *dhtml.HtmlPiece {
	return p.regions.Get("title")
}

func (p *PageTemplate) Main(v any) *PageTemplate {
	p.regions.Add("main", v)
	return p
}

func (p *PageTemplate) GetMain() *dhtml.HtmlPiece {
	return p.regions.Get("main")
}

func (p *PageTemplate) String() string {
	return p.render().String()
}

func (p *PageTemplate) render() (out *dhtml.HtmlPiece) {
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

	container := dhtml.Div().Class("container my-3").
		Append(p.renderHeader())

	// H1 page title
	if !title.IsEmpty() {
		container.Append(dhtml.NewTag("h1").Append(title))
	}

	container.Append(dhtml.Div().Class("main").Append(p.GetMain())).
		Append(p.renderFooter())

	document.Body().Append(container)

	//scripts
	document.Body().
		Append(dhtml.NewTag("script").Attribute("src", "/assets/vendor/bootstrap.bundle.min.js")).
		Append(dhtml.NewTag("script").Attribute("src", "/assets/vendor/vue.global.prod.js")).
		Append(dhtml.NewTag("script").Attribute("src", "/assets/script.min.js"))

	return dhtml.Piece(document)
}

func (p *PageTemplate) renderHeader() (out dhtml.HtmlPiece) {
	var user *model.MtUser
	if v, ok := p.context.Get("User"); ok {
		user = v.(*model.MtUser)
	}

	header := dhtml.Div().Class("border bg-light p-3 mb-3").Attribute("role", "header")

	header_left := dhtml.Div().
		Append(dhtml.Div().Append(dhtml.NewLink("/").Label(app.App.AppName).Class("text-decoration-none"))).
		Append(dhtml.Div().Class("small text-muted").Append("v." + app.App.Version))

	header_right := dhtml.Div().Class("text-end")

	if user != nil {
		header_right.
			Text(user.DisplayName).
			Append(
				dhtml.NewLink("/logout").Label(mtweb.Icon("arrow-right-from-bracket")).
					Class("ms-1").Title("Sign Out"),
			)
	}

	header.Append(mtweb.NewJustifiedLR().L(header_left).R(header_right))

	out.Append(header)
	return out
}

func (p *PageTemplate) renderFooter() (out dhtml.HtmlPiece) {
	out.Append(dhtml.Div().Class("border bg-light p-3 mt-3").Append(
		mtweb.NewJustifiedLR().
			L(fmt.Sprintf("%s v.%s", app.App.AppName, app.App.Version)).
			R(
				dhtml.Div().Class("small text-muted").
					Text("by ").
					Append(dhtml.NewLink("https://www.mito-team.com").Label("MiTo Team").Target("blank")),
			),
	))
	return out
}
