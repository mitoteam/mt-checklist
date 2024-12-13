package app

import (
	"github.com/mitoteam/goappbase"
	gorm "gorm.io/gorm"
)

var (
	App             *goappbase.AppBase
	DefaultSettings *AppSettingsType
	Db              *gorm.DB
)

func init() {
	//default settings (no defaults for now)
	DefaultSettings = &AppSettingsType{}

	//default values for goappbase.AppSettingsBase options
	DefaultSettings.WebserverPort = 15119
}
