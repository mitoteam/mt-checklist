package model

import (
	"reflect"

	"github.com/mitoteam/goapp"
)

type Option struct {
	goapp.BaseModel

	Name  string `gorm:"uniqueIndex"`
	Value string
}

func init() {
	goapp.DbSchema.AddModel(reflect.TypeFor[Option]())
}

// Create and return user
func NewOption(name string) *Option {
	option := &Option{
		Name: name,
	}

	return option
}
