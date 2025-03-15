package yamlpathutils

import (
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestPathToLineNumbers(t *testing.T) {
	// Example YAML content
	yamlContent := `
root:
  child1:
    key1: value1
    key2: value2
  child2:
    key3: value3
`
	var node yaml.Node
	if err := yaml.Unmarshal([]byte(yamlContent), &node); err != nil {
		t.Fatalf("failed to unmarshal YAML: %v", err)
	}

	tests := []struct {
		name     string
		yamlPath string
		expected []int
	}{
		{
			name:     "Match multiple nodes",
			yamlPath: "$.root.child1.*",
			expected: []int{4, 5}, // Lines corresponding to key1 and key2
		},
		{
			name:     "Match single node",
			yamlPath: "$.root.child2.key3",
			expected: []int{7}, // Line corresponding to key3
		},
		{
			name:     "Invalid path",
			yamlPath: "$.root.nonexistent",
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := yamlpath.NewPath(tt.yamlPath)
			if err != nil {
				t.Fatalf("failed to create path: %v", err)
			}

			result := PathToLineNumbers(&node, path)
			if len(result) != len(tt.expected) {
				t.Errorf("unexpected number of results: got %v, want %v", result, tt.expected)
			}

			for i, line := range result {
				if line != tt.expected[i] {
					t.Errorf("unexpected line number at index %d: got %d, want %d", i, line, tt.expected[i])
				}
			}
		})
	}
}

func TestPathToFirstLineNumber(t *testing.T) {
	// Example YAML content
	yamlContent := `
root:
  child1:
    key1: value1
    key2: value2
  child2:
    key3: value3
`
	var node yaml.Node
	if err := yaml.Unmarshal([]byte(yamlContent), &node); err != nil {
		t.Fatalf("failed to unmarshal YAML: %v", err)
	}

	tests := []struct {
		name     string
		yamlPath string
		expected int
	}{
		{
			name:     "Match multiple nodes",
			yamlPath: "$.root.child1.*",
			expected: 4, // First line corresponding to key1
		},
		{
			name:     "Match single node",
			yamlPath: "$.root.child2.key3",
			expected: 7, // Line corresponding to key3
		},
		{
			name:     "Invalid path",
			yamlPath: "$.root.nonexistent",
			expected: -1, // No match
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := yamlpath.NewPath(tt.yamlPath)
			if err != nil {
				t.Fatalf("failed to create path: %v", err)
			}

			result := PathToFirstLineNumber(&node, path)
			if result != tt.expected {
				t.Errorf("unexpected result: got %d, want %d", result, tt.expected)
			}
		})
	}
}
