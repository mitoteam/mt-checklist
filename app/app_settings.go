package app

import (
	"github.com/mitoteam/goapp"
)

// Settings are stored in .settings.yml and not changeable at runtime
type AppSettingsType struct {
	goapp.AppSettingsBase `yaml:",inline"`

	SortOrderStep int64 `yaml:"sort_order_step" yaml_comment:"Default step for sort order items numbers"`

	LdapUrl            string `yaml:"ldap_url" yaml_comment:"LDAPS server URI, for example ldaps://ldap.example.org:636"`
	LdapBindDN         string `yaml:"ldap_bind_dn" yaml_comment:"Service bind DN, for example cn=checklists,ou=checklists,ou=services,dc=example,dc=org"`
	LdapBindPassword   string `yaml:"ldap_bind_password" yaml_comment:"Service bind password; protect this settings file and do not commit real credentials"`
	LdapUserBaseDN     string `yaml:"ldap_user_base_dn" yaml_comment:"User search base DN, for example ou=active,ou=users,dc=example,dc=org"`
	LdapUserFilter     string `yaml:"ldap_user_filter" yaml_comment:"LDAP filter with {username}, for example (&(cn={username})(memberOf=cn=allowed,ou=checklists,ou=services,dc=example,dc=org)); username is escaped before insertion"`
	LdapTimeoutSeconds int    `yaml:"ldap_timeout_seconds" yaml_comment:"LDAP connection and operation timeout in seconds (allowed: 1-60)"`
}

var defaultSettings *AppSettingsType

func init() {
	//default settings (no defaults for now)
	defaultSettings = &AppSettingsType{}

	defaultSettings.SortOrderStep = 10

	//default values for goapp.AppSettingsBase options
	defaultSettings.WebserverPort = 15119

	defaultSettings.LdapTimeoutSeconds = 5
}
