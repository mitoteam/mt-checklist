package model

import (
	"reflect"

	"github.com/mitoteam/goapp"
	"gorm.io/gorm"
)

type ChecklistTemplate struct {
	goapp.BaseModel

	Name          string
	ChecklistName string
}

func init() {
	goapp.DbSchema.AddModel(reflect.TypeFor[ChecklistTemplate]())
}

func (t *ChecklistTemplate) BeforeDelete(tx *gorm.DB) (err error) {
	for _, item := range t.Items() {
		if err := goapp.DeleteObject(item); err != nil {
			return err
		}
	}

	return nil
}

func (t *ChecklistTemplate) Items() []*ChecklistTemplateItem {
	goapp.PreQuery[ChecklistTemplateItem]().Where("checklist_template_id", t.ID)
	return goapp.LoadOL[ChecklistTemplateItem]()
}

func (t *ChecklistTemplate) ItemCount() int64 {
	goapp.PreQuery[ChecklistTemplateItem]().Where("checklist_template_id", t.ID)
	return goapp.CountOL[ChecklistTemplateItem]()
}

type ChecklistTemplateItem struct {
	goapp.BaseModel

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
	goapp.DbSchema.AddModel(reflect.TypeFor[ChecklistTemplateItem]())
}

func (item *ChecklistTemplateItem) GetChecklistTemplate() *ChecklistTemplate {
	if item.ChecklistTemplate == nil {
		item.ChecklistTemplate = goapp.LoadOMust[ChecklistTemplate](item.ChecklistTemplateID)
	}

	return item.ChecklistTemplate
}
