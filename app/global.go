package app

import (
	"github.com/mitoteam/goappbase"
)

var (
	App      *goappbase.AppBase
	Settings *AppSettingsType
)

func init() {
	//default settings (no defaults for now)
	Settings = &AppSettingsType{}

	//default values for goappbase.AppSettingsBase options
	Settings.WebserverPort = 15119
}
