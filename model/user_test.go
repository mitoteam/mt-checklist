package model

import "testing"

func TestNewUserDefaultsToLocalAuthMode(t *testing.T) {
	user := NewUser()
	if user.AuthMode != AuthModeLocal {
		t.Fatalf("AuthMode = %d, want %d", user.AuthMode, AuthModeLocal)
	}
}

func TestAuthModeOptions(t *testing.T) {
	options := AuthModeOptions()
	if len(options) != 2 {
		t.Fatalf("got %d auth mode options, want 2", len(options))
	}
	if options[0].Value != AuthModeLocal || options[0].Label != "Local" {
		t.Fatalf("got auth mode option %+v, want Local", options[0])
	}
	if options[1].Value != AuthModeLDAP || options[1].Label != "LDAP" {
		t.Fatalf("got auth mode option %+v, want LDAP", options[1])
	}
	if !IsValidAuthMode(AuthModeLocal) {
		t.Fatal("Local auth mode should be valid")
	}
	if !IsValidAuthMode(AuthModeLDAP) {
		t.Fatal("LDAP auth mode should be valid")
	}
	if IsValidAuthMode(-1) || IsValidAuthMode(2) {
		t.Fatal("unknown auth modes should be invalid")
	}
}

func TestCheckPasswordOnlyChecksPasswordHash(t *testing.T) {
	user := NewUser()
	user.SetPassword("correct-password")

	if !user.CheckPassword("correct-password") {
		t.Fatal("Local auth mode should accept the correct password")
	}

	user.AuthMode = AuthModeLDAP
	if !user.CheckPassword("correct-password") {
		t.Fatal("CheckPassword should only compare the password hash")
	}
}
