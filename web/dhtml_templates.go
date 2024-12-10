package web

import (
	"github.com/mitoteam/dhtml"
)

type PageTemplate struct {
	// "title" = H1 page title
	// "main" = main content
	regions dhtml.NamedHtmlPieces
}

var pageTemplate *PageTemplate

func GetPageTemplate() *PageTemplate {
	if pageTemplate == nil {
		pageTemplate = &PageTemplate{
			regions: dhtml.NewNamedHtmlPieces(),
		}
	}

	return pageTemplate
}

func (p *PageTemplate) Clear() *PageTemplate {
	p.regions.Clear()
	return p
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

	var head_title = "MT Checklist"

	title := p.regions.Get("title")
	if !title.IsEmpty() {
		head_title = title.String() + " | " + head_title
	}

	document.
		Title(head_title).
		Stylesheet("/assets/vendor/bootstrap.min.css").
		Stylesheet("/assets/vendor/fontawesome.min.css").
		Stylesheet("/assets/vendor/regular.min.css")

	// H1 page title
	if !title.IsEmpty() {
		document.Body().Append(dhtml.NewTag("h1").Append(title))
	}

	// main region
	if r := p.GetMain(); !r.IsEmpty() {
		document.Body().Append(dhtml.Div().Class("main").Append(r))
	}

	return dhtml.Piece(document)
}
