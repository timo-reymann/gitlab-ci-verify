package gitlab_api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGitLabApiClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusTeapot)
		token := request.Header.Get("Private-Token")
		if token == "" {
			t.Fatal("Expected Private-Token header to be present")
		}
	}))

	client := NewGitlabApiClient(server.URL, "token")
	req, err := client.NewRequest("GET", "/", []byte{})
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusTeapot {
		t.Fatalf("Did not contact mock server")
	}
}
