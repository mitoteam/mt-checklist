package app

import (
	"errors"
	"log"

	"github.com/mitoteam/mt-checklist/model"
	gorm "gorm.io/gorm"
)

func checkRootUser() error {
	rootUser := model.User{}

	if err := Db.Limit(1).Find(&rootUser, 1).Error; err != nil {
		log.Println("checkRootUser ERROR: " + err.Error())
		return err
	}

	if rootUser.ID == 0 {
		//root user not found, create one
		rootUser.ID = 1
		rootUser.UserName = "root"
		rootUser.DisplayName = "Root"
		rootUser.SetPassword(App.AppSettings.(*AppSettingsType).InitialRootPassword)

		if err := Db.Create(&rootUser).Error; err != nil {
			log.Println("Db.Create ERROR: " + err.Error())
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

func AuthorizeUser(username, password string) *model.User {
	user := model.User{}

	err := Db.Where(model.User{UserName: username}).First(&user).Error

	if err == nil { //found
		//check password
		if user.CheckPassword(password) {
			return &user
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Query ERROR: " + err.Error())
		return nil
	}

	return nil
}

func GetUser(id int64) *model.User {
	if id == 0 {
		return nil
	}

	user := model.User{}

	err := Db.First(&user, id).Error

	if err == nil { //found
		return &user
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Query ERROR: " + err.Error())
		return nil
	}

	return nil
}
