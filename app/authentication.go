package app

import (
	"crypto/tls"
	"errors"
	"log"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/mitoteam/goapp"
	"github.com/mitoteam/mt-checklist/model"
	"github.com/mitoteam/mttools"
)

var (
	errInvalidCredentials = errors.New("invalid credentials")
	errLDAPUnavailable    = errors.New("LDAP authentication is unavailable")
)

// Checks username and password against the database and/or LDAP server, depending on the user's AuthMode.
// Returns user object if authentication is successful, or nil if user not found or password is wrong and errors
// if there is a problem with the authentication process (for example LDAP server is unavailable)
func AuthenticateUser(username, password string) (*model.User, error) {
	goapp.PreQuery[model.User]().Where("is_active", 1).Where("user_name", username)
	user := goapp.FirstO[model.User]()
	if user == nil {
		return nil, nil // no such user at all
	}

	if err := authenticateUserCredentials(user, password); err != nil {
		if errors.Is(err, errInvalidCredentials) {
			return nil, nil
		}
		return nil, err
	}

	user.LastLogin = mttools.Ptr(time.Now())
	goapp.SaveObject(user)

	return user, nil
}

func authenticateUserCredentials(user *model.User, password string) error {
	switch user.AuthMode {
	case model.AuthModeLocal:
		if user.CheckPassword(password) {
			return nil
		}
		return errInvalidCredentials

	case model.AuthModeLDAP:
		err := authenticateLDAP(user.UserName, password)
		if errors.Is(err, errInvalidCredentials) {
			return errInvalidCredentials
		}
		return err
	default:
		return errors.New("Authentication method is not supported")
	}
}

func validateLDAPSettings(settings *AppSettingsType) error {
	if settings == nil {
		return errors.New("LDAP settings are not configured")
	}

	endpoint, err := url.Parse(settings.LdapUrl)
	if err != nil || endpoint.Scheme != "ldaps" || endpoint.Host == "" || endpoint.User != nil {
		return errors.New("LDAP URL must be a valid ldaps:// URI")
	}
	if strings.TrimSpace(settings.LdapBindDN) == "" ||
		strings.TrimSpace(settings.LdapBindPassword) == "" ||
		strings.TrimSpace(settings.LdapUserBaseDN) == "" {
		return errors.New("LDAP bind credentials and user base DN must be configured")
	}
	if strings.Count(settings.LdapUserFilter, "{username}") != 1 {
		return errors.New("LDAP user filter must contain {username} exactly once")
	}
	if settings.LdapTimeoutSeconds < 1 || settings.LdapTimeoutSeconds > 60 {
		return errors.New("LDAP timeout must be between 1 and 60 seconds")
	}

	return nil
}

func ldapUserFilter(pattern, username string) (string, error) {
	if strings.Count(pattern, "{username}") != 1 {
		return "", errors.New("LDAP user filter must contain {username} exactly once")
	}

	return strings.Replace(pattern, "{username}", ldap.EscapeFilter(username), 1), nil
}

func authenticateLDAP(username, password string) error {
	var settings *AppSettingsType
	var ok bool
	settings, ok = App.AppSettings.(*AppSettingsType)

	if !ok || settings == nil {
		return errors.New("LDAP authentication is not configured")
	}

	if err := validateLDAPSettings(settings); err != nil {
		return err
	}

	filter, err := ldapUserFilter(settings.LdapUserFilter, username)
	if err != nil {
		return err
	}

	timeout := time.Duration(settings.LdapTimeoutSeconds) * time.Second
	connection, err := ldap.DialURL(
		settings.LdapUrl,
		ldap.DialWithDialer(&net.Dialer{Timeout: timeout}),
		ldap.DialWithTLSConfig(&tls.Config{MinVersion: tls.VersionTLS12}),
	)

	if err != nil {
		log.Printf("LDAP connection failed: %v", err)
		return errLDAPUnavailable
	}

	defer connection.Close()
	connection.SetTimeout(timeout)

	if err := connection.Bind(settings.LdapBindDN, settings.LdapBindPassword); err != nil {
		log.Printf("LDAP service bind failed: %v", err)
		return errLDAPUnavailable
	}

	searchRequest := ldap.NewSearchRequest(
		settings.LdapUserBaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		2,
		settings.LdapTimeoutSeconds,
		false,
		filter,
		[]string{"distinguishedName"},
		nil,
	)

	//log.Printf("LDAP search: %v", filter)

	searchResult, err := connection.Search(searchRequest)

	if err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultSizeLimitExceeded) {
			return errInvalidCredentials
		}
		log.Printf("LDAP user search failed: %v", err)
		return errLDAPUnavailable
	}

	if len(searchResult.Entries) != 1 || searchResult.Entries[0].DN == "" {
		return errInvalidCredentials
	}

	if err := connection.Bind(searchResult.Entries[0].DN, password); err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultInvalidCredentials) {
			return errInvalidCredentials
		}
		log.Printf("LDAP user bind failed: %v", err)
		return errLDAPUnavailable
	}

	return nil
}
