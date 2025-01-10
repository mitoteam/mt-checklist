package model

import (
	"reflect"
	"time"

	"github.com/mitoteam/goapp"
	"gorm.io/gorm"
)

type Checklist struct {
	goapp.BaseModel

	Name        string
	Description string

	CreatedByID int64 `gorm:"not null"`
	CreatedBy   *User `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
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

// Count of "done" items and total items count
func (cl *Checklist) DoneItemsCount() (done int64, total int64) {
	for _, item := range cl.Items() {
		if item.DoneAt != nil {
			done++
		}

		total++
	}

	return
}

func (cl *Checklist) GetCreatedBy() *User {
	if cl.CreatedBy == nil {
		cl.CreatedBy = goapp.LoadOMust[User](cl.CreatedByID)
	}

	return cl.CreatedBy
}

func (cl *Checklist) MaxItemSortOrder() int64 {
	sortOrder := int64(0)

	for _, item := range cl.Items() {
		sortOrder = max(sortOrder, item.SortOrder)
	}

	return sortOrder
}

func (cl *Checklist) IsActive() bool {
	done, total := cl.DoneItemsCount()
	return done < total
}

// ====================== checklist items ================================
type ChecklistItem struct {
	goapp.BaseModel

	//fk
	ChecklistID int64      `gorm:"not null;index"`
	Checklist   *Checklist `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`

	Caption   string
	Body      string
	SortOrder int64 `gorm:"not null;index"`
	Weight    int64

	ResponsibleID int64 `gorm:"not null"`
	Responsible   *User `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`

	// user who marked this item as "Done"
	DoneByID    *int64
	DoneBy      *User `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	DoneAt      *time.Time
	DoneComment string
}

func init() {
	goapp.DbSchema.AddModel(reflect.TypeFor[ChecklistItem]())
}

func (ci *ChecklistItem) BeforeDelete(tx *gorm.DB) (err error) {
	tx.Where("checklist_item_id = ?", ci.ID).Delete(&ChecklistItemDependency{})
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

func (item *ChecklistItem) DependenciesCount() int64 {
	goapp.PreQuery[ChecklistItemDependency]().Where("checklist_item_id", item.ID)

	return goapp.CountOL[ChecklistItemDependency]()
}

func (item *ChecklistItem) DependenciesList() []*ChecklistItemDependency {
	goapp.PreQuery[ChecklistItemDependency]().Where("checklist_item_id", item.ID).
		Joins("JOIN checklist_item ci ON ci.ID=checklist_item_id").
		Order("ci.sort_order")

	return goapp.LoadOL[ChecklistItemDependency]()
}

func (item *ChecklistItem) GetDoneBy() *User {
	if item.DoneBy == nil && item.DoneByID != nil {
		item.DoneBy = goapp.LoadOMust[User](*item.DoneByID)
	}

	return item.DoneBy
}

// ====================== checklist item deps ================================

type ChecklistItemDependency struct {
	goapp.DbModel // no ID field

	//this item
	ChecklistItemID int64          `gorm:"not null;index"`
	ChecklistItem   *ChecklistItem `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`

	//depends on this one
	RequireChecklistItemID int64          `gorm:"not null;index"`
	RequireChecklistItem   *ChecklistItem `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
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
