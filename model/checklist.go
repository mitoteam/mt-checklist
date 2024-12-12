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

func LoadChecklist(id any) (cl *Checklist) {
	return goappbase.LoadObject[Checklist](id)
}

func GetChecklistsList() (list []*Checklist) {
	goappbase.DbSchema.Db().Model(&Checklist{}).Find(&list)

	return
}
