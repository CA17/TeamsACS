package app

import (
	"os"
	"testing"
)

func TestBootstrapAdminUsernameDefault(t *testing.T) {
	t.Setenv("TEAMSACS_BOOTSTRAP_ADMIN_USERNAME", "")
	_ = os.Unsetenv("TEAMSACS_BOOTSTRAP_ADMIN_USERNAME")
	if got := bootstrapAdminUsername(); got != defaultBootstrapAdminUsername {
		t.Fatalf("username = %q, want %q", got, defaultBootstrapAdminUsername)
	}
}

func TestBootstrapAdminUsernameFromEnv(t *testing.T) {
	t.Setenv("TEAMSACS_BOOTSTRAP_ADMIN_USERNAME", " ops-admin ")
	if got := bootstrapAdminUsername(); got != "ops-admin" {
		t.Fatalf("username = %q, want ops-admin", got)
	}
}

func TestBootstrapAdminPasswordDefault(t *testing.T) {
	_ = os.Unsetenv("TEAMSACS_BOOTSTRAP_ADMIN_PASSWORD")
	password, fromEnv := bootstrapAdminPassword()
	if fromEnv {
		t.Fatal("expected default password not from env")
	}
	if password != defaultBootstrapAdminPassword {
		t.Fatalf("password = %q, want %q", password, defaultBootstrapAdminPassword)
	}
}

func TestBootstrapAdminPasswordFromEnv(t *testing.T) {
	t.Setenv("TEAMSACS_BOOTSTRAP_ADMIN_PASSWORD", "ChangeMeNow!")
	password, fromEnv := bootstrapAdminPassword()
	if !fromEnv {
		t.Fatal("expected password from env")
	}
	if password != "ChangeMeNow!" {
		t.Fatalf("password = %q, want ChangeMeNow!", password)
	}
}
