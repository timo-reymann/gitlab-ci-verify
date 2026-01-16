package ci_lint_api_proxy

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSendErr(t *testing.T) {
	testCases := []struct {
		name           string
		status         int
		message        string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Not Found error",
			status:         http.StatusNotFound,
			message:        "Resource not found",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"Resource not found"}`,
		},
		{
			name:           "Internal Server Error",
			status:         http.StatusInternalServerError,
			message:        "Something went wrong",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"Something went wrong"}`,
		},
		{
			name:           "Bad Request",
			status:         http.StatusBadRequest,
			message:        "Invalid input",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Invalid input"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			sendErr(recorder, tc.status, tc.message)

			if recorder.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, recorder.Code)
			}

			contentType := recorder.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", contentType)
			}

			body := strings.TrimSpace(recorder.Body.String())
			if body != tc.expectedBody {
				t.Errorf("Expected body %s, got %s", tc.expectedBody, body)
			}
		})
	}
}

func TestErrResponse_JsonSerialization(t *testing.T) {
	testCases := []struct {
		name     string
		response errResponse
		expected string
	}{
		{
			name:     "Simple message",
			response: errResponse{Message: "test error"},
			expected: `{"message":"test error"}`,
		},
		{
			name:     "Empty message",
			response: errResponse{Message: ""},
			expected: `{"message":""}`,
		},
		{
			name:     "Message with special characters",
			response: errResponse{Message: "Error: \"test\" failed"},
			expected: `{"message":"Error: \"test\" failed"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := json.Marshal(tc.response)
			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}

			if string(result) != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, string(result))
			}
		})
	}
}
