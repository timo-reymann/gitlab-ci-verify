package ci_lint_api_proxy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/gitlab/api"
)

func TestAddContextMiddleware(t *testing.T) {
	mockClient := api.NewClient("https://gitlab.example.com", "test-token")
	apiCtx := &ApiContext{
		GitLabClient: mockClient,
	}

	handlerCalled := false
	var receivedCtx *ApiContext

	handler := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		receivedCtx = getApiContext(r)
		w.WriteHeader(http.StatusOK)
	}

	wrappedHandler := addContextMiddleware(handler, apiCtx)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	recorder := httptest.NewRecorder()

	wrappedHandler(recorder, req)

	if !handlerCalled {
		t.Error("Expected handler to be called")
	}

	if receivedCtx == nil {
		t.Fatal("Expected to receive ApiContext, got nil")
	}

	if receivedCtx.GitLabClient != mockClient {
		t.Error("Expected GitLabClient to match the one provided in context")
	}
}

func TestGetApiContext(t *testing.T) {
	mockClient := api.NewClient("https://gitlab.example.com", "test-token")
	apiCtx := &ApiContext{
		GitLabClient: mockClient,
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	reqWithCtx := req.WithContext(context.WithValue(req.Context(), contextKeyApiContext, apiCtx))

	result := getApiContext(reqWithCtx)

	if result == nil {
		t.Fatal("Expected ApiContext, got nil")
	}

	if result.GitLabClient != mockClient {
		t.Error("Expected GitLabClient to match")
	}
}
