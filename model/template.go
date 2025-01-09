package model

import (
	"reflect"

	"github.com/mitoteam/goapp"
	"gorm.io/gorm"
)

type Template struct {
	goapp.BaseModel

	Name                 string // template name
	ChecklistName        string // initial name to create checklist
	ChecklistDescription string // initial description to create checklist
}

func init() {
	goapp.DbSchema.AddModel(reflect.TypeFor[Template]())
}

func (t *Template) BeforeDelete(tx *gorm.DB) (err error) {
	for _, item := range t.Items() {
		if err := goapp.DeleteObject(item); err != nil {
			return err
		}
	}

	return nil
}

func (t *Template) Items() []*TemplateItem {
	goapp.PreQuery[TemplateItem]().Where("template_id", t.ID).Order("sort_order")
	return goapp.LoadOL[TemplateItem]()
}

func (t *Template) ItemCount() int64 {
	goapp.PreQuery[TemplateItem]().Where("template_id", t.ID)
	return goapp.CountOL[TemplateItem]()
}

func (t *Template) MaxItemSortOrder() int64 {
	sortOrder := int64(0)

	for _, item := range t.Items() {
		sortOrder = max(sortOrder, item.SortOrder)
	}

	return sortOrder
}

// ===================== template items ========================
type TemplateItem struct {
	goapp.BaseModel

	//fk
	TemplateID int64     `gorm:"not null;index"`
	Template   *Template `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`

	Caption   string
	Body      string
	SortOrder int64 `gorm:"not null;index"`
	Weight    int64

	ResponsibleID int64 `gorm:"not null;index"`
	Responsible   *User `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

func init() {
	goapp.DbSchema.AddModel(reflect.TypeFor[TemplateItem]())
}

func (ti *TemplateItem) BeforeDelete(tx *gorm.DB) (err error) {
	tx.Where("template_item_id = ?", ti.ID).Delete(&TemplateItemDependency{})
	return nil
}

func (item *TemplateItem) GetTemplate() *Template {
	if item.Template == nil {
		item.Template = goapp.LoadOMust[Template](item.TemplateID)
	}

	return item.Template
}

func (item *TemplateItem) GetResponsible() *User {
	if item.Responsible == nil {
		item.Responsible = goapp.LoadOMust[User](item.ResponsibleID)
	}

	return item.Responsible
}

func (item *TemplateItem) DependenciesCount() int64 {
	goapp.PreQuery[TemplateItemDependency]().Where("template_item_id", item.ID)

	return goapp.CountOL[TemplateItemDependency]()
}

func (item *TemplateItem) DependenciesList() []*TemplateItemDependency {
	goapp.PreQuery[TemplateItemDependency]().Where("template_item_id", item.ID).
		Joins("JOIN template_item ti ON ti.ID=template_item_id").
		Order("ti.sort_order")

	return goapp.LoadOL[TemplateItemDependency]()
}

// ======================= item dependencies ============================
type TemplateItemDependency struct {
	goapp.DbModel // no ID field

	//this item
	TemplateItemID int64         `gorm:"not null;index"`
	TemplateItem   *TemplateItem `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`

	//depends on this one
	RequireTemplateItemID int64         `gorm:"not null"`
	RequireTemplateItem   *TemplateItem `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

func init() {
	goapp.DbSchema.AddModel(reflect.TypeFor[TemplateItemDependency]())
}

func (item *TemplateItemDependency) GetTemplateItem() *TemplateItem {
	if item.TemplateItem == nil {
		item.TemplateItem = goapp.LoadOMust[TemplateItem](item.TemplateItemID)
	}

	return item.TemplateItem
}

func (item *TemplateItemDependency) GetRequireTemplateItem() *TemplateItem {
	if item.RequireTemplateItem == nil {
		item.RequireTemplateItem = goapp.LoadOMust[TemplateItem](item.RequireTemplateItemID)
	}

	return item.RequireTemplateItem
}
