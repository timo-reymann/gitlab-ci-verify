package includes

import (
	"gopkg.in/yaml.v3"
	"testing"
)

func TestComponentInclude_Type(t *testing.T) {
	componentInclude := &ComponentInclude{}
	expectedType := "component"
	if componentInclude.Type() != expectedType {
		t.Fatalf("Type() = %v, want %v", componentInclude.Type(), expectedType)
	}
}

func TestComponentInclude_Equals(t *testing.T) {
	tests := []struct {
		name     string
		include1 *ComponentInclude
		include2 Include
		expected bool
	}{
		{
			name: "equal includes",
			include1: &ComponentInclude{
				Component: "my-component",
			},
			include2: &ComponentInclude{
				Component: "my-component",
			},
			expected: true,
		},
		{
			name: "different components",
			include1: &ComponentInclude{
				Component: "my-component-1",
			},
			include2: &ComponentInclude{
				Component: "my-component-2",
			},
			expected: false,
		},
		{
			name: "different types",
			include1: &ComponentInclude{
				Component: "my-component",
			},
			include2: &LocalInclude{
				Path: "test.yml",
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

func TestNewComponentInclude(t *testing.T) {
	node := &yaml.Node{}
	component := "my-component"
	componentInclude := NewComponentInclude(node, component)

	if componentInclude.Component != component {
		t.Fatalf("NewComponentInclude().Component = %v, want %v", componentInclude.Component, component)
	}

	if componentInclude.Node() != node {
		t.Fatalf("NewComponentInclude().Node() = %v, want %v", componentInclude.Node(), node)
	}
}
