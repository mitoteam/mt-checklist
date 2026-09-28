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
	if len(options) != 1 {
		t.Fatalf("got %d auth mode options, want 1", len(options))
	}
	if options[0].Value != AuthModeLocal || options[0].Label != "Local" {
		t.Fatalf("got auth mode option %+v, want Local", options[0])
	}
	if !IsValidAuthMode(AuthModeLocal) {
		t.Fatal("Local auth mode should be valid")
	}
	if IsValidAuthMode(-1) || IsValidAuthMode(1) {
		t.Fatal("unknown auth modes should be invalid")
	}
}

func TestCheckPasswordRequiresLocalAuthMode(t *testing.T) {
	user := NewUser()
	user.SetPassword("correct-password")

	if !user.CheckPassword("correct-password") {
		t.Fatal("Local auth mode should accept the correct password")
	}

	user.AuthMode = 99
	if user.CheckPassword("correct-password") {
		t.Fatal("unknown auth mode should reject password authentication")
	}
}
