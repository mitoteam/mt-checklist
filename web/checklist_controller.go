package web

import (
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mtweb"
)

type ChecklistController struct {
	mbr.ControllerBase
}

var ClCtl *ChecklistController

func init() {
	ClCtl = &ChecklistController{}
	ClCtl.With(AuthMiddleware)
}

func (c *RootController) Checklist() mbr.Route {
	return mbr.Route{PathPattern: "/checklist/{checklist_id}", ChildController: ClCtl}
}

func (c *ChecklistController) Checklist() mbr.Route {
	return mbr.Route{
		PathPattern: "/",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			cl := model.LoadChecklist(p.ctx.Request().PathValue("checklist_id"))

			p.Title(cl.Name)

			descriptionOut := dhtml.Div().Class("mb-3")

			if cl.Description != "" {
				descriptionOut.Append(dhtml.Div().Class("text-prewrap").Append(cl.Description))
			}

			p.Main(descriptionOut)

			if len(cl.Items()) > 0 {

				cardList := mtweb.NewCardList()

				for _, item := range cl.Items() {
					bodyOut := dhtml.NewHtmlPiece()

					if item.Body != "" {
						bodyOut.Append(dhtml.Div().Class("text-prewrap").Append(item.Body))
					}

					card := mtweb.NewCard().
						Header(item.Caption).
						Body(bodyOut)

					cardList.Add(card)
				}

				p.Main(cardList)
			} else {
				p.Main(dhtml.EmptyLabel("no items added"))
			}

			return nil
		}),
	}
}
