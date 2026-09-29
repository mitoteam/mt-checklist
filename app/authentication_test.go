package app

import (
	"testing"
)

func TestValidateLDAPSettings(t *testing.T) {
	validSettings := &AppSettingsType{
		LdapUrl:            "ldaps://ldap.example.org:636",
		LdapBindDN:         "cn=checklist,ou=service,dc=example,dc=org",
		LdapBindPassword:   "secret",
		LdapUserBaseDN:     "ou=users,dc=example,dc=org",
		LdapUserFilter:     "(uid={username})",
		LdapTimeoutSeconds: 5,
	}
	if err := validateLDAPSettings(validSettings); err != nil {
		t.Fatalf("valid LDAP settings rejected: %v", err)
	}

	tests := []struct {
		name     string
		settings *AppSettingsType
	}{
		{"non-LDAPS URL", func() *AppSettingsType { c := *validSettings; c.LdapUrl = "ldap://ldap.example.org"; return &c }()},
		{"missing filter placeholder", func() *AppSettingsType { c := *validSettings; c.LdapUserFilter = "(uid=alice)"; return &c }()},
		{"invalid timeout", func() *AppSettingsType { c := *validSettings; c.LdapTimeoutSeconds = 0; return &c }()},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := validateLDAPSettings(test.settings); err == nil {
				t.Fatal("expected invalid LDAP settings to be rejected")
			}
		})
	}
}

func TestLDAPUserFilterEscapesUsername(t *testing.T) {
	filter, err := ldapUserFilter("(uid={username})", "a*)(uid=*)")
	if err != nil {
		t.Fatalf("ldapUserFilter returned error: %v", err)
	}
	if filter != "(uid=a\\2a\\29\\28uid=\\2a\\29)" {
		t.Fatalf("filter = %q", filter)
	}
}
