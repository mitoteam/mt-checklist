package model

import (
	"reflect"

	"github.com/mitoteam/goappbase"
)

type Checklist struct {
	goappbase.BaseModel

	Name     string
	IsActive bool
	Body     string `gorm:"type:varchar(1000)"`
}

type ChecklistItem struct {
	goappbase.BaseModel

	Name string
	Body string `gorm:"type:varchar(1000)"`
}

func init() {
	goappbase.DbSchema.AddModel(reflect.TypeFor[Checklist]())
}

func LoadChecklist(id any) *Checklist {
	return goappbase.LoadO[Checklist](id)
}
