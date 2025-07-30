package api

import (
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/netrc"
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/vault"
	"os"
	"strings"
)

type TokenSource interface {
	Lookup(baseUrl TokenSourceLookupHints) (string, error)
}

type TokenSourceLookupHints struct {
	ExistingToken string
	BaseUrl       string
}

type EnvVarTokenSource struct{}

func (e EnvVarTokenSource) Lookup(_ TokenSourceLookupHints) (string, error) {
	return os.Getenv("GITLAB_TOKEN"), nil
}

type NetRcTokenSource struct{}

func (n NetRcTokenSource) Lookup(t TokenSourceLookupHints) (string, error) {
	userNetrc, err := netrc.ReadUserNetrc()
	if err != nil {
		return "", err
	}

	credentials, err := netrc.GetCredentials(userNetrc, t.BaseUrl)
	if err != nil {
		return "", err
	}

	return credentials.Password, nil
}

type VaultTokenSource struct{}

func (v VaultTokenSource) Lookup(t TokenSourceLookupHints) (string, error) {
	if !strings.HasPrefix(t.ExistingToken, "vault://") && !strings.ContainsRune(t.ExistingToken, '#') {
		return "", nil
	}

	spec := t.ExistingToken[8:len(t.ExistingToken)]
	specParts := strings.SplitN(spec, "#", 2)
	if len(specParts) != 2 {
		return "", fmt.Errorf("malformed spec '%s', should be in format 'vault://<field>#<path>'", t.ExistingToken)
	}

	return vault.GetSecretWithDefaultConfiguration(specParts[0], specParts[1])
}
