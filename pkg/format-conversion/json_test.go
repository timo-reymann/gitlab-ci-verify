package format_conversion

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestParseJson(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected map[string]any
		wantErr  bool
	}{
		{
			name:     "valid JSON",
			input:    []byte(`{"key1": "value1", "key2": 123}`),
			expected: map[string]any{"key1": "value1", "key2": 123},
			wantErr:  false,
		},
		{
			name:     "invalid JSON",
			input:    []byte(`{invalid json}`),
			expected: map[string]any{},
			wantErr:  true,
		},
		{
			name:     "empty JSON",
			input:    []byte(`{}`),
			expected: map[string]any{},
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := ParseJson(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("ParseJson(%v): expected error %v, got %v", tc.input, tc.wantErr, err)
				return
			}
			if res == nil {
				t.Fatalf("Expected result to be deserialized")
			}
		})
	}
}

func TestToJson(t *testing.T) {
	testCases := []struct {
		name     string
		input    map[string]any
		expected []byte
		isError  bool
	}{
		{
			name:     "Empty map",
			input:    map[string]any{},
			expected: []byte("{}"),
			isError:  false,
		},
		{
			name:     "Simple map",
			input:    map[string]any{"key": "value"},
			expected: []byte(`{"key":"value"}`),
			isError:  false,
		},
		{
			name:     "Invalid JSON value",
			input:    map[string]any{"key": func() {}},
			expected: nil,
			isError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ToJson(tc.input)

			if tc.isError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if !cmp.Equal(result, tc.expected) {
					t.Errorf("Expected: %s, Got: %s", string(tc.expected), string(result))
				}
			}
		})
	}
}
