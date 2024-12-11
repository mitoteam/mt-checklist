package app

import (
	"errors"
	"log"

	"github.com/mitoteam/mt-checklist/model"
	gorm "gorm.io/gorm"
)

func GetChecklist(id int64) (cl *model.MtChecklist) {
	if id == 0 {
		return
	}

	err := Db.First(&cl, id).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Query ERROR: " + err.Error())
		return nil
	}

	return
}

func GetChecklistsList() (list []*model.MtChecklist) {
	Db.Model(&model.MtChecklist{}).Find(&list)

	return
}
