package app

import (
	"log"

	"github.com/mitoteam/goappbase"
	"github.com/mitoteam/mt-checklist/model"
	gorm "gorm.io/gorm"
)

var (
	App      *goappbase.AppBase
	Settings *AppSettingsType
	Db       *gorm.DB
)

func init() {
	//default settings (no defaults for now)
	Settings = &AppSettingsType{}

	//default values for goappbase.AppSettingsBase options
	Settings.WebserverPort = 15119
}

func DoPreRun() error {
	var err error

	// open database and migrate schema
	if Db, err = goappbase.DbSchema.Open(); err != nil {
		return err
	}

	//check if root user exists
	rootUser := model.MtUser{}

	if err := Db.Limit(1).Find(&rootUser, 1).Error; err != nil {
		log.Println("FirstOrInit ERROR: " + err.Error())
		return err
	}

	if rootUser.ID == 0 {
		//root user not found, create one
		rootUser.ID = 1
		rootUser.UserName = "root"
		rootUser.DisplayName = "Root"
		rootUser.SetPassword(App.AppSettings.(*AppSettingsType).InitialRootPassword)

		if err := Db.Create(&rootUser).Error; err != nil {
			log.Println("BotDb.Create ERROR: " + err.Error())
			return err
		}

		log.Printf(
			"Root user created with initial password '%s'. PLEASE CHANGE IT AS SOON AS POSSIBLE!\n",
			App.AppSettings.(*AppSettingsType).InitialRootPassword,
		)
	}

	//log.Printf("%+v\n", rootUser)

	return nil //no errors
}

func DoPostRun() error {
	goappbase.DbSchema.Close()

	return nil //no errors
}
