package model

import "github.com/mitoteam/goappbase"

type MtChecklist struct {
	goappbase.BaseModel

	Name string
}

func init() {
	goappbase.DbSchema.AddModel(&MtChecklist{})
}
