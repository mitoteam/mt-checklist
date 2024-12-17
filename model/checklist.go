package model

import (
	"reflect"

	"github.com/mitoteam/goapp"
)

type Checklist struct {
	goapp.BaseModel

	Name        string
	IsActive    bool
	Description string
}

type ChecklistItem struct {
	goapp.BaseModel

	Name string
	Body string
}

func init() {
	goapp.DbSchema.AddModel(reflect.TypeFor[Checklist]())
}

func LoadChecklist(id any) *Checklist {
	return goapp.LoadO[Checklist](id)
}
