package api

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestGitLabApiClient(t *testing.T) {
	server := mockServerWithoutBodyAndVerifier(http.StatusTeapot, func(request *http.Request) {
		token := request.Header.Get("Private-Token")
		if token == "" {
			t.Fatal("Expected Private-Token header to be present")
		}
	})

	client := NewClient(server.URL, "token")
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

func TestClient_LintCiYaml(t *testing.T) {
	testCases := []struct {
		name           string
		projectSlug    string
		response       string
		expectedResult *CiLintResult
	}{
		{
			name:        "valid",
			projectSlug: "org/group/project",
			response:    `{ "valid": true, "mergedYaml": "merged_yaml", "errors": [], "warnings": [] }`,
			expectedResult: &CiLintResult{
				Valid:      true,
				MergedYaml: "merged_yaml",
				Errors:     []string{},
				Warnings:   []string{},
			},
		},
		{
			name:        "warnings",
			projectSlug: "org/group/project",
			response:    `{ "valid": false, "mergedYaml": "merged_yaml", "errors": [], "warnings": ["warning text"] }`,
			expectedResult: &CiLintResult{
				Valid:      false,
				MergedYaml: "merged_yaml",
				Errors:     []string{},
				Warnings:   []string{"warning text"},
			},
		},
		{
			name:        "errors",
			projectSlug: "org/group/project",
			response:    `{ "valid": false, "mergedYaml": "merged_yaml", "errors": ["error text"], "warnings": [] }`,
			expectedResult: &CiLintResult{
				Valid:      false,
				MergedYaml: "merged_yaml",
				Errors:     []string{"error text"},
				Warnings:   []string{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := mockServerWithBodyAndVerifier(http.StatusTeapot, []byte(tc.response), func(request *http.Request) {
				if !strings.Contains(request.URL.RawPath, "%2F") {
					log.Fatal("project parameter is not url escaped properly")
				}
			})
			defer server.Close()

			client := NewClient(server.URL, "token")
			result, err := client.LintCiYaml(context.Background(), tc.projectSlug, []byte(""))
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(result, tc.expectedResult) {
				t.Fatalf("Expected result to be %v, but got %v", tc.expectedResult, result)
			}
		})
	}

}
