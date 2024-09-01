package format_conversion

import (
	"testing"
)

func TestParseYaml(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected map[any]any
		wantErr  bool
	}{
		{
			name: "valid YAML",
			input: []byte(`key1: value1
key2: 123`),
			expected: map[any]any{"key1": "value1", "key2": 123},
			wantErr:  false,
		},
		{
			name:     "invalid YAML",
			input:    []byte(`invalid yaml`),
			expected: map[any]any{},
			wantErr:  true,
		},
		{
			name:     "empty YAML",
			input:    []byte(""),
			expected: map[any]any{},
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseYaml(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("ParseYaml(%v): expected error %v, got %v", tc.input, tc.wantErr, err)
				return
			}
		})
	}
}