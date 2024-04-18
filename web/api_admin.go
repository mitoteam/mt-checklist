package web

import (
	"github.com/mitoteam/goappbase"
)

type apiChecklist struct {
	Name string `json:"name"`
}

func api_AdminChecklists(r *goappbase.ApiRequest) error {
	if !apiCheckSession(r) {
		return nil
	}

	r.SetOutData("checklists", []apiChecklist{{Name: "First"}, {Name: "Second"}})

	return nil
}
