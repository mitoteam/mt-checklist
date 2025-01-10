package app

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mt-checklist/model"
)

func MySqlMigrate() {
	db, err := sql.Open("mysql", "mt_checklists:X1tI2lneW@/checklists_mito_team_com")
	if err != nil {
		panic(err)
	}

	//check connection
	if err = db.Ping(); err != nil {
		panic(err.Error())
	}

	log.Println("Migration: Mysql PING OK")

	var checklistRows, itemRows, depRows *sql.Rows
	var (
		checklistId, itemId int64
		checklistName       string
		datetimeStr         string
	)

	// ======================= migrate templates ==========================
	if checklistRows, err = db.Query("select ID, Name, Created from mtc_Checklists WHERE IsTemplate = 1"); err != nil {
		panic(err)
	}
	defer checklistRows.Close()

	for checklistRows.Next() {
		if err = checklistRows.Scan(&checklistId, &checklistName, &datetimeStr); err != nil {
			panic(err)
		}

		itemRows, err = db.Query(
			"SELECT ID, ShortName, Name, Sequence, Weight FROM mtc_Items WHERE ChecklistID=?", checklistId,
		)
		if err != nil {
			panic(err)
		}
		defer itemRows.Close()

		depRows, err = db.Query(
			"SELECT MasterID, SlaveID FROM mtc_Dependency WHERE SlaveID IN (SELECT ID FROM mtc_Items WHERE ChecklistID=?)", checklistId,
		)
		if err != nil {
			panic(err)
		}
		defer depRows.Close()

		itemMap := make(map[int64]int64)

		goapp.Transaction(func() error {
			template := &model.Template{
				Name:          checklistName,
				ChecklistName: checklistName,
			}

			template.CreatedAt, _ = time.Parse(time.DateTime, datetimeStr)

			goapp.SaveObject(template)

			for itemRows.Next() {
				templateItem := &model.TemplateItem{
					TemplateID:    template.ID,
					ResponsibleID: 1, //root
				}

				err = itemRows.Scan(
					&itemId, &templateItem.Caption, &templateItem.Body, &templateItem.SortOrder, &templateItem.Weight,
				)
				if err != nil {
					panic(err)
				}

				goapp.SaveObject(templateItem)

				itemMap[itemId] = templateItem.ID
			}

			for depRows.Next() {
				var masterID, slaveID int64

				err = depRows.Scan(&masterID, &slaveID)
				if err != nil {
					panic(err)
				}

				//log.Printf("DEP: %d -> %d\n", masterID, slaveID)
				dep := &model.TemplateItemDependency{
					TemplateItemID:        itemMap[slaveID],
					RequireTemplateItemID: itemMap[masterID],
				}

				goapp.CreateObject(dep)
			}

			return nil
		})

		log.Printf("Template '%s' DONE", checklistName)
	}

	// ======================= migrate checklists ==========================

	if checklistRows, err = db.Query("select ID, Name, Created from mtc_Checklists WHERE IsTemplate = 0"); err != nil {
		panic(err)
	}
	defer checklistRows.Close()

	for checklistRows.Next() {
		if err = checklistRows.Scan(&checklistId, &checklistName, &datetimeStr); err != nil {
			panic(err)
		}

		itemRows, err = db.Query(
			"SELECT ID, ShortName, Name, Sequence, Weight, IFNULL(CheckedDate, ''), IFNULL(Comment, '') FROM mtc_Items WHERE ChecklistID=?",
			checklistId,
		)
		if err != nil {
			panic(err)
		}
		defer itemRows.Close()

		depRows, err = db.Query(
			"SELECT MasterID, SlaveID FROM mtc_Dependency WHERE SlaveID IN (SELECT ID FROM mtc_Items WHERE ChecklistID=?)", checklistId,
		)
		if err != nil {
			panic(err)
		}
		defer depRows.Close()

		itemMap := make(map[int64]int64)

		goapp.Transaction(func() error {
			checklist := &model.Checklist{
				Name:        checklistName,
				Description: "imported from old version",
				CreatedByID: 1, // root
			}

			checklist.CreatedAt, _ = time.Parse(time.DateTime, datetimeStr)

			goapp.SaveObject(checklist)

			for itemRows.Next() {
				checklistItem := &model.ChecklistItem{
					ChecklistID:   checklist.ID,
					ResponsibleID: 1, //alway root since we are not importing users
				}

				err = itemRows.Scan(
					&itemId, &checklistItem.Caption, &checklistItem.Body, &checklistItem.SortOrder,
					&checklistItem.Weight, &datetimeStr, &checklistItem.DoneComment,
				)
				if err != nil {
					panic(err)
				}

				if datetimeStr != "" {
					timeValue, _ := time.Parse(time.DateTime, datetimeStr)
					checklistItem.DoneAt = &timeValue
					checklistItem.DoneByID = 1 //root
				}

				goapp.SaveObject(checklistItem)

				itemMap[itemId] = checklistItem.ID
			}

			for depRows.Next() {
				var masterID, slaveID int64

				err = depRows.Scan(&masterID, &slaveID)
				if err != nil {
					panic(err)
				}

				//log.Printf("DEP: %d -> %d\n", masterID, slaveID)
				dep := &model.ChecklistItemDependency{
					ChecklistItemID:        itemMap[slaveID],
					RequireChecklistItemID: itemMap[masterID],
				}

				goapp.CreateObject(dep)
			}

			return nil
		})

		log.Printf("Checklist '%s' MIGRATION DONE", checklistName)
	}

	if err = checklistRows.Err(); err != nil {
		panic(err)
	}

	defer db.Close()
}
