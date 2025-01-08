package web

import (
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mt-checklist/model"
)

func renderChecklistItemBody(item *model.ChecklistItem) (out dhtml.HtmlPiece) {
	if item.Body == "" {
		out.Append(dhtml.Div().Append(dhtml.EmptyLabel("no description")))
	} else {
		out.Append(dhtml.Div().Append(MdEngine.ToDhtml(item.Body)))
	}

	out.Append(dhtml.RenderValue("Responsible", item.GetResponsible().GetDisplayName()).Class("mt-3"))

	//dependencies
	cellOut := dhtml.Div()
	if item.RequiredItemsCount() > 0 {
		depsList := dhtml.NewUnorderedList().Class("mb-0")

		for _, dep := range item.RequiredItems() {
			depsList.AppendItem(dhtml.NewListItem().Append(dep.GetRequireChecklistItem().Caption))
		}

		cellOut.Append(dhtml.Div().Class("fs-5").Append("Depends"))
		cellOut.Append(depsList)

		out.Append(cellOut)
	}

	return out
}

func createChecklistFromTemplate(template *model.Template) *model.Checklist {
	checklist := &model.Checklist{}
	checklist.IsActive = true
	checklist.Name = template.ChecklistName
	checklist.Description = template.ChecklistDescription

	goapp.Transaction(func() error {
		goapp.SaveObject(checklist) //we need an ID to create items

		for _, templateItem := range template.Items() {
			checklistItem := &model.ChecklistItem{
				ChecklistID:   checklist.ID,
				Caption:       templateItem.Caption,
				Body:          templateItem.Body,
				Weight:        templateItem.Weight,
				SortOrder:     templateItem.SortOrder,
				ResponsibleID: templateItem.ResponsibleID,
			}

			goapp.SaveObject(checklistItem)
		}

		return nil
	})

	return checklist
}
