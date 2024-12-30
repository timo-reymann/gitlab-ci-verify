package api

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"os"
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
			response:    `{ "valid": true, "merged_yaml": "merged_yaml", "errors": [], "warnings": [] }`,
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
			response:    `{ "valid": false, "merged_yaml": "merged_yaml", "errors": [], "warnings": ["warning text"] }`,
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
			response:    `{ "valid": false, "merged_yaml": "merged_yaml", "errors": ["error text"], "warnings": [] }`,
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
					t.Fatal("project parameter is not url escaped properly")
				}
			})
			defer server.Close()

			client := NewClient(server.URL, "token")
			result, err := client.LintCiYaml(context.TODO(), tc.projectSlug, []byte(""))
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(result, tc.expectedResult) {
				t.Fatalf("Expected result to be %v, but got %v", tc.expectedResult, result)
			}
		})
	}

}

func TestNewClientWithMultiTokenSources(t *testing.T) {
	testCases := []struct {
		name          string
		token         string
		netrcFile     string
		gitlabToken   string
		expectedToken string
		expectedError bool
	}{
		{
			name:          "Token parameter is non-empty",
			token:         "my-token",
			expectedToken: "my-token",
		},
		{
			name:          "Netrc contains a host entry with a password",
			netrcFile:     "test.netrc",
			expectedToken: "my-password",
		},
		{
			name:          "Environment variable GITLAB_TOKEN is set",
			gitlabToken:   "my-env-token",
			expectedToken: "my-env-token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.netrcFile != "" {
				err := os.WriteFile(tc.netrcFile, []byte("machine gitlab.example.com login user password my-password"), 0600)
				if err != nil {
					t.Fatalf("Failed to create netrc file: %v", err)
				}
				defer os.Remove(tc.netrcFile)

				// Set the NETRC environment variable to point to the temporary file
				_ = os.Setenv("NETRC", tc.netrcFile)
			}

			if tc.gitlabToken != "" {
				_ = os.Setenv("GITLAB_TOKEN", tc.gitlabToken)
			}

			client := NewClientWithMultiTokenSources("https://gitlab.example.com", tc.token)
			if client.token != tc.expectedToken {
				t.Errorf("Expected token to be '%s', but got '%s'", tc.expectedToken, client.token)
			}

			if tc.netrcFile != "" {
				_ = os.Remove(tc.netrcFile)
				_ = os.Unsetenv("NETRC")
			}
			if tc.gitlabToken != "" {
				_ = os.Unsetenv("GITLAB_TOKEN")
			}
		})
	}
}
