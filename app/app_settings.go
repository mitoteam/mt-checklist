package app

import (
	"github.com/mitoteam/goapp"
)

type AppSettingsType struct {
	goapp.AppSettingsBase `yaml:",inline"`

	ExampleOption string `yaml:"example_option" yaml_comment:"TODO: remove"`
}

var defaultSettings *AppSettingsType

func init() {
	//default settings (no defaults for now)
	defaultSettings = &AppSettingsType{}

	//default values for goapp.AppSettingsBase options
	defaultSettings.WebserverPort = 15119
}
