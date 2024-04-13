package model

import "github.com/mitoteam/goappbase"

type ChUser struct {
	goappbase.BaseModel

	Name string
}

func init() {
	goappbase.DbSchema.AddModel(&ChUser{})
}
