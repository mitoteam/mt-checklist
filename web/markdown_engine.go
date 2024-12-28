package web

import (
	"bytes"

	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/mttools"
	"github.com/mitoteam/mtweb"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type markdownEngine struct {
	mttools.Lockable //make it goroutine safe
	engine           goldmark.Markdown
}

var MdEngine *markdownEngine

func init() {
	MdEngine = &markdownEngine{
		engine: goldmark.New(
			goldmark.WithExtensions(extension.Linkify),
		),
	}
}

func (e *markdownEngine) ToDhtml(input string) *dhtml.HtmlPiece {
	var buf bytes.Buffer

	//goroutine safety
	e.Lock()
	err := e.engine.Convert([]byte(input), &buf)
	e.Unlock()

	if err != nil {
		return mtweb.RenderErrorf("Error in commonmark: %s", err.Error())
	}

	return dhtml.Piece(dhtml.UnsafeText(buf.String()))
}
