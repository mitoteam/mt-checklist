package app

import (
	"errors"
	"log"

	"github.com/mitoteam/mt-checklist/model"
	gorm "gorm.io/gorm"
)

func GetChecklist(id int64) *model.MtChecklist {
	if id == 0 {
		return nil
	}

	o := model.MtChecklist{}

	err := Db.First(&o, id).Error

	if err == nil { //found
		return &o
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Query ERROR: " + err.Error())
		return nil
	}

	return nil
}

func GetChecklistsList() (list []*model.MtChecklist) {
	rows, _ := Db.Model(&model.MtChecklist{}).Rows()
	defer rows.Close()

	o := model.MtChecklist{}
	for rows.Next() {
		Db.ScanRows(rows, &o)

		new_o := o
		list = append(list, &new_o)
	}

	return list
}
