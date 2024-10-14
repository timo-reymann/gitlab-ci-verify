package vault

import (
	"errors"
	"os"
)

// GetSecretWithDefaultConfiguration using the token configured by the CLI and the
// environment variable VAULT_ADDR containing the base URL for the vault API
func GetSecretWithDefaultConfiguration(path string, field string) (string, error) {
	token, err := LookupToken()
	if err != nil {
		return "", err
	}

	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		return "", errors.New("environment variable VAULT_ADDR not set, which is required for fetching secret")
	}

	return getSecret(vaultAddr, token, path, field)
}
