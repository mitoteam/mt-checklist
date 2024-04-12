package app

import (
	"github.com/mitoteam/goappbase"
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
