package web

import (
	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mtweb"
)

func renderChecklistItemBody(item *model.ChecklistItem) (out dhtml.HtmlPiece) {
	if item.DoneByID != nil {
		out.Append(dhtml.Div().Append(mtweb.Icon(iconUser).Label(item.GetDoneBy().GetDisplayName())))
		out.Append(mtweb.NewTimestamp(*item.DoneAt).Icon("square-check"))
	} else {
		out.Append(dhtml.Div().Append(mtweb.Icon(iconUser).Label(item.GetResponsible().GetDisplayName())))
	}

	if item.Body != "" {
		out.Append(dhtml.Div().Class("mt-3").Append(MdEngine.ToDhtml(item.Body)))
	}

	//dependencies
	cellOut := dhtml.Div()
	if len(item.GetUnresolvedDepItemList()) > 0 {
		depsList := dhtml.NewUnorderedList().Class("mb-0")

		for _, depItem := range item.GetUnresolvedDepItemList() {
			depsList.AppendItem(dhtml.NewListItem().Append(depItem.Caption))
		}

		cellOut.Append(
			dhtml.Div().Append(mtweb.Icon(mtweb.IconNameNo).Label("Unresolved:").ElementClass("fw-bold text-danger")),
			depsList,
		)

		out.Append(cellOut)
	}

	return out
}

func createChecklistFromTemplate(template *model.Template, user *model.User) *model.Checklist {
	checklist := &model.Checklist{}
	checklist.Name = template.ChecklistName
	checklist.Description = template.ChecklistDescription
	checklist.CreatedByID = user.ID

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
