package format_conversion

import (
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
