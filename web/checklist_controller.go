package web

import (
	"time"

	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mbr"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mttools"
	"github.com/mitoteam/mtweb"
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
func (c *ChecklistController) ViewChecklist() mbr.Route {
	return mbr.Route{
		PathPattern: "/",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			cl := model.LoadChecklist(p.ctx.Request().PathValue("checklist_id"))

			p.Title(cl.Name)

			/// INFO and DESCRIPTION
			descriptionOut := dhtml.Div().Class("mb-3")
			descriptionOut.Append(
				dhtml.RenderValue(
					"Created by",
					dhtml.Span().Append(
						cl.GetCreatedBy().GetDisplayName(),
						mtweb.NewTimestamp(cl.CreatedAt).SmallMuted().Class("d-inline ms-3"),
					),
				),
			)

			descriptionOut.Append(dhtml.RenderValue("Progress", cl.GetProgress()))

			if cl.Description != "" {
				descriptionOut.Append(dhtml.Div().Class("mt-3").Append(MdEngine.ToDhtml(cl.Description)))
			}

			p.Main(descriptionOut)

			/// ITEMS
			if len(cl.Items()) > 0 {
				cardList := dhtmlbs.NewCardList().Class("row-cols-1", "row-cols-lg-2", "row-cols-xxl-3")

				for _, item := range orderedChecklistItems(cl, p.User()) {
					status := item.GetStatus(p.User())

					var captionOut dhtml.HtmlPiece
					if status == model.ITEM_STATUS_RED {
						captionOut.Append(mtweb.Icon(mtweb.IconNameNo).Class("me-1"))
					} else if status == model.ITEM_STATUS_YELLOW {
						captionOut.Append(mtweb.Icon("triangle-exclamation").Class("me-1"))
					}

					captionOut.Append(dhtml.Span().Class("fw-bold").Append(item.Caption))

					header := dhtmlbs.NewJustifiedLR().L(
						dhtml.NewTag("b").Append(captionOut),
					)

					if item.CanDone() {
						header.R(
							mtweb.NewIconBtnR(
								"check", "",
								ChecklistCtl.ChecklistItemDone, "checklist_id", cl.ID, "item_id", item.ID,
							).Class("btn-success py-1").Title("Mark as done").Confirm("Mark item as done?"),
						)
					}

					card := dhtmlbs.NewCard().
						Header(header).
						Body(renderChecklistItemBody(item, p.User()))

					if status == model.ITEM_STATUS_GREEN {
						card.Class("border-success").BodyClass("bg-success-subtle")
					} else if status == model.ITEM_STATUS_RED {
						card.Class("border-danger", "opacity-50").BodyClass("bg-danger-subtle").HeaderClass("text-bg-danger")
					} else if status == model.ITEM_STATUS_YELLOW {
						card.Class("border-warning").BodyClass("bg-warning-subtle").HeaderClass("text-bg-warning")
					}

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

// checklist page
func (c *ChecklistController) ChecklistItemDone() mbr.Route {
	return mbr.Route{
		PathPattern: "/item/{item_id}/done",
		HandleF: PageBuilderRouteHandler(func(p *PageBuilder) any {
			cl := model.LoadChecklist(p.ctx.Request().PathValue("checklist_id"))
			item := goapp.LoadOrCreateO[model.ChecklistItem](p.ctx.Request().PathValue("item_id"))

			mttools.AssertEqual(item.ChecklistID, cl.ID)

			item.DoneAt = mttools.Ptr(time.Now())
			item.DoneByID = &p.User().ID //current user

			goapp.SaveObject(item)

			p.RedirectRoute(ChecklistCtl.ViewChecklist, "checklist_id", cl.ID)

			return nil
		}),
	}
}
