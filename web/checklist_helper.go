package web

import (
	"slices"

	"github.com/mitoteam/dhtml"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mtweb"
)

func renderChecklistItemBody(item *model.ChecklistItem, user *model.User) (out dhtml.HtmlPiece) {
	if item.DoneByID != nil {
		out.Append(dhtml.Div().Append(mtweb.Icon(iconUser).Label(item.GetDoneBy().GetDisplayName())))
		out.Append(mtweb.NewTimestamp(*item.DoneAt).Icon("square-check"))
	} else {
		iconName := iconUser
		if item.GetStatus(user) == model.ITEM_STATUS_YELLOW {
			iconName = "triangle-exclamation"
		}

		respOut := dhtml.Piece(mtweb.Icon(iconName).Label(item.GetResponsible().GetDisplayName()))

		out.Append(dhtml.Div().Append(respOut))
	}

	if item.Body != "" {
		out.Append(dhtml.Div().Class("mt-3", "mt-no-last-p-margin").Append(MdEngine.ToDhtml(item.Body)))
	}

	//unresolved dependencies
	if len(item.GetUnresolvedDepItemList()) > 0 {
		depsOut := dhtml.Div().Class("mt-3")

		depsList := dhtml.NewUnorderedList().Class("mb-0")

		for _, depItem := range item.GetUnresolvedDepItemList() {
			depsList.AppendItem(dhtml.NewListItem().Append(depItem.Caption))
		}

		depsOut.Append(
			dhtml.Div().Append(mtweb.Icon(mtweb.IconNameNo).Label("Unresolved:").ElementClass("fw-bold text-danger")),
			depsList,
		)

		out.Append(depsOut)
	}

	return out
}

func orderedChecklistItems(cl *model.Checklist, user *model.User) (list []*model.ChecklistItem) {
	list = slices.Clone(cl.Items())
	slices.SortStableFunc(list, func(a, b *model.ChecklistItem) int {
		return a.GetStatus(user) - b.GetStatus(user)
	})

	return list
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
