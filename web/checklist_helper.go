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
	if item.DependenciesCount() > 0 {
		depsList := dhtml.NewUnorderedList().Class("mb-0")

		for _, dep := range item.DependenciesList() {
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
	checklist.Name = template.ChecklistName
	checklist.Description = template.ChecklistDescription

	goapp.Transaction(func() error {
		goapp.SaveObject(checklist) //we need an ID to create items

		// templateItemID => checklistItem
		itemMap := make(map[int64]*model.ChecklistItem, template.ItemCount())

		//create items
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

			itemMap[templateItem.ID] = checklistItem
		}

		//create dependencies
		for _, templateItem := range template.Items() {
			checklistItem := itemMap[templateItem.ID]

			for _, requiredTemplateItem := range templateItem.DependenciesList() {
				dep := &model.ChecklistItemDependency{
					ChecklistItemID:        checklistItem.ID,
					RequireChecklistItemID: itemMap[requiredTemplateItem.RequireTemplateItemID].ID,
				}

				goapp.CreateObject(dep)
			}
		}

		return nil
	})

	return checklist
}
