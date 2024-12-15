package model

import (
	"reflect"

	"github.com/mitoteam/goappbase"
)

type ChecklistTemplate struct {
	goappbase.BaseModel

	Name          string
	ChecklistName string

	Items []ChecklistTemplateItem
}

type ChecklistTemplateItem struct {
	goappbase.BaseModel

	//fk
	ChecklistTemplateID int64

	Caption   string
	Body      string
	SortOrder int
	Weight    int

	ResponsibleID uint
	Responsible   User
}

func init() {
	goappbase.DbSchema.AddModel(reflect.TypeFor[ChecklistTemplate]())
	goappbase.DbSchema.AddModel(reflect.TypeFor[ChecklistTemplateItem]())
}
