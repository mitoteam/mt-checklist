package app

import (
	"log"

	"github.com/mitoteam/mt-checklist/model"
)

func checkRootUser() error {
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
	return nil
}
