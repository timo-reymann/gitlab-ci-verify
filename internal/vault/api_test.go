package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSecret(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(`{ "data": {"data": { "password": "password-val" }} }`))
	}))
	secret, err := getSecret(srv.URL, "token", "path/to/secret", "password")
	if err != nil {
		t.Fatal(err)
	}
	if secret != "password-val" {
		t.Fatal("Failed to query")
	}
	defer srv.Close()
}
