package web

import (
	"github.com/mitoteam/dhtmlbs"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mttools"
)

func NewUserSelect(name string) *dhtmlbs.SelectFormControlElement {
	selectControl := dhtmlbs.NewSelect(name)

	//add users
	goapp.PreQuery[model.User]().Where("is_active", 1).Order("user_name")
	for _, user := range goapp.LoadOL[model.User]() {
		selectControl.Option(mttools.AnyToString(user.ID), user.GetDisplayName())
	}

	return selectControl
}
