package ci_lint_api_proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/gitlab/api"
)

func TestProjectLintHandler_MethodNotAllowed(t *testing.T) {
	methods := []string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/v4/projects/test/ci/lint", nil)
			recorder := httptest.NewRecorder()

			projectLintHandler(recorder, req)

			if recorder.Code != http.StatusMethodNotAllowed {
				t.Errorf("Expected status %d for method %s, got %d", http.StatusMethodNotAllowed, method, recorder.Code)
			}

			var response errResponse
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Message != "Only POST requests are allowed" {
				t.Errorf("Expected message 'Only POST requests are allowed', got '%s'", response.Message)
			}
		})
	}
}

func TestProjectLintHandler_Success(t *testing.T) {
	gitlabResponse := `{"valid":true,"merged_yaml":"test","errors":[],"warnings":[]}`
	gitlabServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request to GitLab, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "/api/v4/projects/") {
			t.Errorf("Expected path to contain /api/v4/projects/, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(gitlabResponse))
	}))
	defer gitlabServer.Close()

	mockClient := api.NewClient(gitlabServer.URL, "test-token")
	apiCtx := &ApiContext{
		GitLabClient: mockClient,
	}

	requestBody := `{"content":"stages:\n  - test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v4/projects/myproject/ci/lint", strings.NewReader(requestBody))
	req.SetPathValue("project", "myproject")
	req = req.WithContext(context.WithValue(req.Context(), contextKeyApiContext, apiCtx))
	recorder := httptest.NewRecorder()

	projectLintHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	contentType := recorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	body := strings.TrimSpace(recorder.Body.String())
	if body != gitlabResponse {
		t.Errorf("Expected body to match GitLab response.\nExpected: %s\nGot: %s", gitlabResponse, body)
	}
}

func TestProjectLintHandler_GitLabError(t *testing.T) {
	gitlabServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer gitlabServer.Close()

	mockClient := api.NewClient(gitlabServer.URL, "test-token")
	apiCtx := &ApiContext{
		GitLabClient: mockClient,
	}

	requestBody := `{"content":"invalid"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v4/projects/myproject/ci/lint", strings.NewReader(requestBody))
	req.SetPathValue("project", "myproject")
	req = req.WithContext(context.WithValue(req.Context(), contextKeyApiContext, apiCtx))
	recorder := httptest.NewRecorder()

	projectLintHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Logf("Handler returned status %d (GitLab response proxied)", recorder.Code)
	}
}

func TestProjectLintHandler_WithPathParameter(t *testing.T) {
	var receivedPath string
	gitlabServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"valid":true}`))
	}))
	defer gitlabServer.Close()

	mockClient := api.NewClient(gitlabServer.URL, "test-token")
	apiCtx := &ApiContext{
		GitLabClient: mockClient,
	}

	testCases := []struct {
		name           string
		projectPath    string
		expectedInPath string
	}{
		{
			name:           "Simple project",
			projectPath:    "myproject",
			expectedInPath: "/api/v4/projects/myproject/ci/lint",
		},
		{
			name:           "Numeric project ID",
			projectPath:    "12345",
			expectedInPath: "/api/v4/projects/12345/ci/lint",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			receivedPath = ""
			req := httptest.NewRequest(http.MethodPost, "/api/v4/projects/"+tc.projectPath+"/ci/lint", strings.NewReader(`{}`))
			req.SetPathValue("project", tc.projectPath)
			req = req.WithContext(context.WithValue(req.Context(), contextKeyApiContext, apiCtx))
			recorder := httptest.NewRecorder()

			projectLintHandler(recorder, req)

			if receivedPath != tc.expectedInPath {
				t.Errorf("Expected GitLab to receive path %s, got %s", tc.expectedInPath, receivedPath)
			}
		})
	}
}

func TestProjectLintHandler_RequestBodyForwarded(t *testing.T) {
	var receivedBody []byte
	gitlabServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		receivedBody, err = io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed to read request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"valid":true}`))
	}))
	defer gitlabServer.Close()

	mockClient := api.NewClient(gitlabServer.URL, "test-token")
	apiCtx := &ApiContext{
		GitLabClient: mockClient,
	}

	requestBody := `{"content":"stages:\n  - build\n  - test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v4/projects/myproject/ci/lint", strings.NewReader(requestBody))
	req.SetPathValue("project", "myproject")
	req = req.WithContext(context.WithValue(req.Context(), contextKeyApiContext, apiCtx))
	recorder := httptest.NewRecorder()

	projectLintHandler(recorder, req)

	if !bytes.Equal(receivedBody, []byte(requestBody)) {
		t.Errorf("Request body not forwarded correctly.\nExpected: %s\nGot: %s", requestBody, string(receivedBody))
	}
}
