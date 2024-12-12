package model

import (
	"reflect"

	"github.com/mitoteam/goappbase"
)

type Checklist struct {
	goappbase.BaseModel

	Name string
	Body string `gorm:"type:varchar(1000)"`
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
	return goappbase.LoadObject[Checklist](id)
}

func GetChecklistsList() []*Checklist {
	return goappbase.LoadObjectList[Checklist]()
}
