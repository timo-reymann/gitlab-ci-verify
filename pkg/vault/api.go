package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type vaultResponse struct {
	Data struct {
		Data map[string]string `json:"data"`
	} `json:"data"`
}

func getSecret(addr string, token string, path string, field string) (string, error) {
	pathParts := strings.SplitN(path, "/", 2)
	engine, secret := pathParts[0], pathParts[1]

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/%s/data/%s", strings.TrimSuffix(addr, "/"), engine, secret), bytes.NewBuffer([]byte{}))
	if err != nil {
		return "", err
	}

	req.Header.Add("X-Vault-Token", token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	content, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid HTTP status: %d - %s", res.StatusCode, string(content))
	}

	var decoded vaultResponse
	err = json.Unmarshal(content, &decoded)
	if err != nil {
		return "", err
	}

	val, ok := decoded.Data.Data[field]
	if !ok {
		return "", fmt.Errorf("no field '%s' in secret for path '%s'", field, path)
	}

	return val, nil
}
