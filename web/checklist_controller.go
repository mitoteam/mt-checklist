package web

import (
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-checklist/model"
)

type ChecklistController struct {
	mbr.ControllerBase
}

var ChecklistCtl *ChecklistController

func init() {
	ChecklistCtl = &ChecklistController{}
	ChecklistCtl.With(AuthMiddleware)
}

// base route for ChecklistController
func (c *RootController) Checklist() mbr.Route {
	return mbr.Route{PathPattern: "/checklist/{checklist_id}", ChildController: ChecklistCtl}
}

// checklist page
func (c *ChecklistController) Checklist() mbr.Route {
	return mbr.Route{
		PathPattern: "/",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			cl := model.LoadChecklist(p.ctx.Request().PathValue("checklist_id"))

			p.Title(cl.Name)

			descriptionOut := dhtml.Div().Class("mb-3")

			if cl.Description != "" {
				descriptionOut.Append(dhtml.Div().Append(MdEngine.ToDhtml(cl.Description)))
			}

			p.Main(descriptionOut)

			if len(cl.Items()) > 0 {

				cardList := dhtmlbs.NewCardList()

				for _, item := range cl.Items() {
					card := dhtmlbs.NewCard().
						Header(item.Caption).
						Body(renderChecklistItemBody(item))

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
