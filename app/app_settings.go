package app

import (
	"github.com/mitoteam/goapp"
)

// Settings are stored in .settings.yml and not changeable at runtime
type AppSettingsType struct {
	goapp.AppSettingsBase `yaml:",inline"`

	SortOrderStep int64 `yaml:"sort_order_step" yaml_comment:"Default step for sort order items numbers"`
}

var defaultSettings *AppSettingsType

func init() {
	//default settings (no defaults for now)
	defaultSettings = &AppSettingsType{}

	defaultSettings.SortOrderStep = 10

	//default values for goapp.AppSettingsBase options
	defaultSettings.WebserverPort = 15119
}
