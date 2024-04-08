package app

import (
	"github.com/mitoteam/goappbase"
)

type AppSettingsType struct {
	goappbase.AppSettingsBase `yaml:",inline"`

	BotToken string `yaml:"bot_token" yaml_comment:"Bot authorization token"`
}
