package api

import (
	"context"
	"testing"
)

func TestNewRequest(t *testing.T) {
	testCases := []struct {
		name               string
		method             string
		baseUrl            string
		path               string
		token              string
		payload            []byte
		expectedPath       string
		expectedAuthHeader string
	}{
		{
			name:               "Valid request",
			method:             "GET",
			baseUrl:            "https://api.example.com",
			path:               "/users",
			token:              "your_token",
			payload:            nil,
			expectedPath:       "/users",
			expectedAuthHeader: "your_token",
		},
		{
			name:               "Empty path",
			method:             "POST",
			baseUrl:            "https://api.example.com",
			path:               "/",
			token:              "your_token",
			payload:            []byte("{\"name\":\"John\"}"),
			expectedPath:       "/",
			expectedAuthHeader: "your_token",
		},
		{
			name:               "Path starting with slash",
			method:             "PUT",
			baseUrl:            "https://api.example.com",
			path:               "/projects/1/",
			token:              "your_token",
			payload:            nil,
			expectedPath:       "/projects/1/",
			expectedAuthHeader: "your_token",
		},
		{
			name:               "Invalid method",
			method:             "INVALID_METHOD",
			baseUrl:            "https://api.example.com",
			path:               "/projects/2",
			token:              "your_token",
			payload:            nil,
			expectedPath:       "/projects/2",
			expectedAuthHeader: "your_token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := NewRequest(context.Background(), tc.method, tc.baseUrl, tc.path, tc.token, tc.payload)
			if err != nil {
				if tc.expectedPath == "" {
					t.Logf("Expected error: %v", err)
				} else {
					t.Fatalf("NewRequest failed: %v", err)
				}
				return
			}

			if req.URL.Path != tc.expectedPath {
				t.Errorf("Path is incorrect: expected %q, got %q", tc.expectedPath, req.URL.Path)
			}

			authHeader := req.Header.Get("PRIVATE-TOKEN")
			if authHeader != tc.expectedAuthHeader {
				t.Errorf("Authorization header is incorrect: expected %q, got %q", tc.expectedAuthHeader, authHeader)
			}
		})
	}
}
