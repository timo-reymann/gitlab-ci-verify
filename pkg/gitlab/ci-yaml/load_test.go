package ci_yaml

import (
	"reflect"
	"testing"
)

func TestLoadRaw(t *testing.T) {
	testCases := []struct {
		name          string
		path          string
		expectedMap   map[string]any
		expectedError bool
	}{
		{
			name:          "Valid YAML file",
			path:          "test_data/valid.yaml",
			expectedMap:   map[string]any{"key1": "value1", "key2": 123},
			expectedError: false,
		},
		{
			name:          "Invalid YAML file",
			path:          "test_data/invalid.yaml",
			expectedMap:   nil,
			expectedError: true,
		},
		{
			name:          "Non-existent file",
			path:          "test_data/non_existent.yaml",
			expectedMap:   nil,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := LoadRaw(tc.path)
			if err != nil {
				if !tc.expectedError {
					t.Errorf("Unexpected error: %v", err)
				}
				return
			}

			if !reflect.DeepEqual(result, tc.expectedMap) {
				t.Errorf("Result is incorrect: expected %v, got %v", tc.expectedMap, result)
			}
		})
	}
}
