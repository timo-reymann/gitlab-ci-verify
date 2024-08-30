package api

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestCheckStatus(t *testing.T) {
	testCases := []struct {
		name        string
		statusCode  int
		expectedErr error
	}{
		{
			name:        "Valid response - 200 OK",
			statusCode:  http.StatusOK,
			expectedErr: nil,
		},
		{
			name:        "Forbidden - 403",
			statusCode:  http.StatusForbidden,
			expectedErr: ErrInvalidAuthentication,
		},
		{
			name:        "Not Found - 404",
			statusCode:  http.StatusNotFound,
			expectedErr: ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := &Response{
				Response: &http.Response{
					StatusCode: tc.statusCode,
				},
			}
			err := resp.CheckStatus()
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUnmarshalJson(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		jsonData   string
		wantErr    error
	}{
		{
			name:       "Valid JSON",
			statusCode: http.StatusOK,
			jsonData:   `{"key": "value"}`,
			wantErr:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := &Response{
				Response: &http.Response{
					StatusCode: tc.statusCode,
					Body:       io.NopCloser(strings.NewReader(tc.jsonData)),
				},
			}

			var val interface{}
			err := resp.UnmarshalJson(&val)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Expected error %v, got %v", tc.wantErr, err)
			}
		})
	}
}
