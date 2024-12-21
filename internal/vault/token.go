package vault

import (
	"errors"
	"os"
	"path"
)

func lookupCliToken() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if dir == "" {
		return "", errors.New("empty home directory")
	}

	token, err := os.ReadFile(path.Join(dir, ".vault-token"))
	if err != nil {
		return "", err
	}

	return string(token), nil
}

func lookupEnvToken() string {
	return os.Getenv("VAULT_TOKEN")
}

// LookupToken to authenticate for API calls
// Checks the env var VAULT_TOKEN first and if that is not set,
// tries to look up ~/.vault-token, which is written by Vault CLI
func LookupToken() (string, error) {
	token := lookupEnvToken()
	if token != "" {
		return token, nil
	}

	return lookupCliToken()
}
