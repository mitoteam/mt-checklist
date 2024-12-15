package model

import (
	"reflect"

	"github.com/mitoteam/goappbase"
)

type ChecklistTemplate struct {
	goappbase.BaseModel

	Name          string
	ChecklistName string
}

type ChecklistTemplateItem struct {
	goappbase.BaseModel

	Name string
}

func init() {
	goappbase.DbSchema.AddModel(reflect.TypeFor[ChecklistTemplate]())
	goappbase.DbSchema.AddModel(reflect.TypeFor[ChecklistTemplateItem]())
}
