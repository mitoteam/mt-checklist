package model

import (
	"reflect"

	"github.com/mitoteam/goappbase"
)

type ChecklistTemplate struct {
	goappbase.BaseModel

	Name          string
	ChecklistName string

	SortOrder3 string `gorm:"index"`
}

func (t *ChecklistTemplate) Items() []*ChecklistTemplateItem {
	goappbase.PreQuery[ChecklistTemplateItem]().Where("checklist_template_id", t.ID)
	return goappbase.LoadOL[ChecklistTemplateItem]()
}

func (t *ChecklistTemplate) ItemCount() int64 {
	goappbase.PreQuery[ChecklistTemplateItem]().Where("checklist_template_id", t.ID)
	return goappbase.CountOL[ChecklistTemplateItem]()
}

type ChecklistTemplateItem struct {
	goappbase.BaseModel

	//fk
	ChecklistTemplateID int64 //`gorm:"not null,index,constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	ChecklistTemplate   *ChecklistTemplate

	Caption   string
	Body      string
	SortOrder int `gorm:"not null,index"`
	Weight    int

	ResponsibleID int64
}

func (item *ChecklistTemplateItem) GetChecklistTemplate() *ChecklistTemplate {
	if item.ChecklistTemplate == nil {
		item.ChecklistTemplate = goappbase.LoadOMust[ChecklistTemplate](item.ChecklistTemplateID)
	}

	return item.ChecklistTemplate
}

func init() {
	goappbase.DbSchema.AddModel(reflect.TypeFor[ChecklistTemplate]())
	goappbase.DbSchema.AddModel(reflect.TypeFor[ChecklistTemplateItem]())
}
