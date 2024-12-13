package app

import (
	"github.com/mitoteam/goappbase"
)

type AppSettingsType struct {
	goappbase.AppSettingsBase `yaml:",inline"`

	ExampleOption string `yaml:"example_option" yaml_comment:"TODO: remove"`
}

var defaultSettings *AppSettingsType

func init() {
	//default settings (no defaults for now)
	defaultSettings = &AppSettingsType{}

	//default values for goappbase.AppSettingsBase options
	defaultSettings.WebserverPort = 15119
}
