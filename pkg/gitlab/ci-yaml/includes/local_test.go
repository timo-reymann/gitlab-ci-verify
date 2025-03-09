package includes

import (
	"gopkg.in/yaml.v3"
	"testing"
)

func TestLocalInclude_Type(t *testing.T) {
	localInclude := &LocalInclude{}
	expectedType := "local"
	if localInclude.Type() != expectedType {
		t.Fatalf("Type() = %v, want %v", localInclude.Type(), expectedType)
	}
}

func TestLocalInclude_Equals(t *testing.T) {
	tests := []struct {
		name     string
		include1 *LocalInclude
		include2 Include
		expected bool
	}{
		{
			name: "equal includes",
			include1: &LocalInclude{
				Path: "test.yml",
			},
			include2: &LocalInclude{
				Path: "test.yml",
			},
			expected: true,
		},
		{
			name: "different paths",
			include1: &LocalInclude{
				Path: "test1.yml",
			},
			include2: &LocalInclude{
				Path: "test2.yml",
			},
			expected: false,
		},
		{
			name: "different types",
			include1: &LocalInclude{
				Path: "test.yml",
			},
			include2: &ProjectInclude{
				Project: "my-group/my-project",
				Files:   []string{"test.yml"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.include1.Equals(tt.include2)
			if result != tt.expected {
				t.Fatalf("Equals() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewLocalInclude(t *testing.T) {
	node := &yaml.Node{}
	path := "test.yml"
	localInclude := NewLocalInclude(node, path)

	if localInclude.Path != path {
		t.Fatalf("NewLocalInclude().Path = %v, want %v", localInclude.Path, path)
	}

	if localInclude.Node() != node {
		t.Fatalf("NewLocalInclude().Node() = %v, want %v", localInclude.Node(), node)
	}
}
