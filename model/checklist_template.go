package model

import (
	"reflect"

	"github.com/mitoteam/goappbase"
	"gorm.io/gorm"
)

type ChecklistTemplate struct {
	goappbase.BaseModel

	Name          string
	ChecklistName string
}

func init() {
	goappbase.DbSchema.AddModel(reflect.TypeFor[ChecklistTemplate]())
}

func (t *ChecklistTemplate) BeforeDelete(tx *gorm.DB) (err error) {
	for _, item := range t.Items() {
		if err := goappbase.DeleteObject(item); err != nil {
			return err
		}
	}

	return nil
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

func init() {
	goappbase.DbSchema.AddModel(reflect.TypeFor[ChecklistTemplateItem]())
}

func (item *ChecklistTemplateItem) GetChecklistTemplate() *ChecklistTemplate {
	if item.ChecklistTemplate == nil {
		item.ChecklistTemplate = goappbase.LoadOMust[ChecklistTemplate](item.ChecklistTemplateID)
	}

	return item.ChecklistTemplate
}
