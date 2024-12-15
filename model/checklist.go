package model

import (
	"reflect"

	"github.com/mitoteam/goappbase"
)

type Checklist struct {
	goappbase.BaseModel

	Name        string
	IsActive    bool
	Description string
}

type ChecklistItem struct {
	goappbase.BaseModel

	Name string
	Body string
}

func init() {
	goappbase.DbSchema.AddModel(reflect.TypeFor[Checklist]())
}

func LoadChecklist(id any) *Checklist {
	return goappbase.LoadO[Checklist](id)
}
