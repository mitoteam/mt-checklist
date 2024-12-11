package web

import (
	"github.com/mitoteam/mtweb"
)

func PageAdminChecklists(p *PageBuilder) bool {
	pageOut := mtweb.Icon("vial")

	p.Main(pageOut)
	return true
}
