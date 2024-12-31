package model

import (
	"reflect"

	"github.com/mitoteam/goapp"
	"gorm.io/gorm"
)

type Checklist struct {
	goapp.BaseModel

	Name        string
	Description string
	IsActive    bool
}

func init() {
	goapp.DbSchema.AddModel(reflect.TypeFor[Checklist]())
}

func LoadChecklist(id any) *Checklist {
	return goapp.LoadOMust[Checklist](id)
}

func (t *Checklist) BeforeDelete(tx *gorm.DB) (err error) {
	for _, item := range t.Items() {
		if err := goapp.DeleteObject(item); err != nil {
			return err
		}
	}

	return nil
}

func (t *Checklist) Items() []*ChecklistItem {
	goapp.PreQuery[ChecklistItem]().Where("checklist_id", t.ID).Order("sort_order")
	return goapp.LoadOL[ChecklistItem]()
}

func (t *Checklist) ItemCount() int64 {
	goapp.PreQuery[ChecklistItem]().Where("checklist_id", t.ID)
	return goapp.CountOL[ChecklistItem]()
}

type ChecklistItem struct {
	goapp.BaseModel

	//fk
	ChecklistID int64 //`gorm:"not null,index,constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	Checklist   *Checklist

	Caption   string
	Body      string
	SortOrder int64 `gorm:"not null,index"`
	Weight    int64

	ResponsibleID int64 //`gorm:"not null,index,constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	Responsible   *User
}

func init() {
	goapp.DbSchema.AddModel(reflect.TypeFor[ChecklistItem]())
}

func (item *ChecklistItem) GetChecklist() *Checklist {
	if item.Checklist == nil {
		item.Checklist = goapp.LoadOMust[Checklist](item.ChecklistID)
	}

	return item.Checklist
}
