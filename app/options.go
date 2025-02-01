package app

import (
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mt-checklist/model"
)

// Options can be changed in runtime
type AppOptionsType struct {
}

var Options *AppOptionsType

func init() {
	Options = &AppOptionsType{}
}

func (o *AppOptionsType) getString(name string, defValue string) string {
	goapp.PreQuery[model.Option]().Where("name", name)
	option := goapp.FirstO[model.Option]()

	if option == nil {
		return defValue
	}

	return option.Value
}

func (o *AppOptionsType) setString(name, value string) {
	goapp.PreQuery[model.Option]().Where("name", name)
	option := goapp.FirstO[model.Option]()

	if option == nil {
		option = model.NewOption(name)
	}

	option.Value = value

	goapp.SaveObject(option)
}

func (o *AppOptionsType) SiteName() string {
	return o.getString("site_name", App.AppName)
}

func (o *AppOptionsType) SetSiteName(value string) {
	o.setString("site_name", value)
}

func (o *AppOptionsType) SiteMotto() string {
	return o.getString("site_motto", "v"+App.Version)
}

func (o *AppOptionsType) SetSiteMotto(value string) {
	o.setString("site_motto", value)
}
