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

func (cl *Checklist) BeforeDelete(tx *gorm.DB) (err error) {
	for _, item := range cl.Items() {
		if err := goapp.DeleteObject(item); err != nil {
			return err
		}
	}

	return nil
}

func (cl *Checklist) Items() []*ChecklistItem {
	goapp.PreQuery[ChecklistItem]().Where("checklist_id", cl.ID).Order("sort_order")
	return goapp.LoadOL[ChecklistItem]()
}

func (cl *Checklist) ItemCount() int64 {
	goapp.PreQuery[ChecklistItem]().Where("checklist_id", cl.ID)
	return goapp.CountOL[ChecklistItem]()
}

// ====================== checklist items ================================
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

func (ci *ChecklistItem) BeforeDelete(tx *gorm.DB) (err error) {
	for _, item := range ci.RequiredItems() {
		if err := goapp.DeleteObject(item); err != nil {
			return err
		}
	}

	return nil
}

func (item *ChecklistItem) GetChecklist() *Checklist {
	if item.Checklist == nil {
		item.Checklist = goapp.LoadOMust[Checklist](item.ChecklistID)
	}

	return item.Checklist
}

func (item *ChecklistItem) GetResponsible() *User {
	if item.Responsible == nil {
		item.Responsible = goapp.LoadOMust[User](item.ResponsibleID)
	}

	return item.Responsible
}

func (item *ChecklistItem) RequiredItemsCount() int64 {
	goapp.PreQuery[ChecklistItemDependency]().Where("checklist_item_id", item.ID)

	return goapp.CountOL[ChecklistItemDependency]()
}

func (item *ChecklistItem) RequiredItems() []*ChecklistItemDependency {
	goapp.PreQuery[ChecklistItemDependency]().Where("checklist_item_id", item.ID).
		Joins("JOIN checklist_item ci ON ci.ID=checklist_item_id").
		Order("ci.sort_order")

	return goapp.LoadOL[ChecklistItemDependency]()
}

// ====================== checklist item deps ================================

type ChecklistItemDependency struct {
	goapp.DbModel // no ID field

	//this item
	ChecklistItemID int64 //`gorm:"not null,index,constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	ChecklistItem   *ChecklistItem

	//depends on this one
	RequireChecklistItemID int64 //`gorm:"not null,index,constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	RequireChecklistItem   *ChecklistItem
}

func init() {
	goapp.DbSchema.AddModel(reflect.TypeFor[ChecklistItemDependency]())
}

func (item *ChecklistItemDependency) GetChecklistItem() *ChecklistItem {
	if item.ChecklistItem == nil {
		item.ChecklistItem = goapp.LoadOMust[ChecklistItem](item.ChecklistItemID)
	}

	return item.ChecklistItem
}

func (item *ChecklistItemDependency) GetRequireChecklistItem() *ChecklistItem {
	if item.RequireChecklistItem == nil {
		item.RequireChecklistItem = goapp.LoadOMust[ChecklistItem](item.RequireChecklistItemID)
	}

	return item.RequireChecklistItem
}
