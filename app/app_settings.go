package app

import (
	"github.com/mitoteam/goappbase"
)

type AppSettingsType struct {
	goappbase.AppSettingsBase `yaml:",inline"`

	ExampleOption string `yaml:"example_option" yaml_comment:"TODO: remove"`
}
