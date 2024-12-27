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
	goapp.PreQuery[TemplateItem]().Where("template_id", t.ID)
	return goapp.LoadOL[TemplateItem]()
}

func (t *Template) ItemCount() int64 {
	goapp.PreQuery[TemplateItem]().Where("template_id", t.ID)
	return goapp.CountOL[TemplateItem]()
}

type TemplateItem struct {
	goapp.BaseModel

	//fk
	TemplateID int64 //`gorm:"not null,index,constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	Template   *Template

	Caption   string
	Body      string
	SortOrder int `gorm:"not null,index"`
	Weight    int

	ResponsibleID int64
}

func init() {
	goapp.DbSchema.AddModel(reflect.TypeFor[TemplateItem]())
}

func (item *TemplateItem) GetTemplate() *Template {
	if item.Template == nil {
		item.Template = goapp.LoadOMust[Template](item.TemplateID)
	}

	return item.Template
}
