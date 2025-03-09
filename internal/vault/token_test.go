package vault

import (
	"os"
	"testing"
)

func TestLookupCliToken_HomeDirError(t *testing.T) {
	originalHomeDir := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHomeDir)

	os.Setenv("HOME", "")
	_, err := lookupCliToken()
	if err == nil {
		t.Fatal("lookupCliToken() error = nil, want non-nil")
	}
}

func TestLookupEnvToken(t *testing.T) {
	expectedToken := "env-token"
	os.Setenv("VAULT_TOKEN", expectedToken)
	defer os.Unsetenv("VAULT_TOKEN")

	token := lookupEnvToken()
	if token != expectedToken {
		t.Fatalf("lookupEnvToken() = %v, want %v", token, expectedToken)
	}
}

func TestLookupToken_EnvToken(t *testing.T) {
	expectedToken := "env-token"
	os.Setenv("VAULT_TOKEN", expectedToken)
	defer os.Unsetenv("VAULT_TOKEN")

	token, err := LookupToken()
	if err != nil {
		t.Fatalf("LookupToken() error = %v, want nil", err)
	}
	if token != expectedToken {
		t.Fatalf("LookupToken() = %v, want %v", token, expectedToken)
	}
}

func TestLookupToken_NoToken(t *testing.T) {
	os.Unsetenv("VAULT_TOKEN")
	originalHomeDir := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHomeDir)

	os.Setenv("HOME", "")
	_, err := LookupToken()
	if err == nil {
		t.Fatal("LookupToken() error = nil, want non-nil")
	}
}
