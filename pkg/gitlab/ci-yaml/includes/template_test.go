package includes

import (
	"gopkg.in/yaml.v3"
	"testing"
)

func TestTemplateInclude_Type(t *testing.T) {
	templateInclude := &TemplateInclude{}
	expectedType := "template"
	if templateInclude.Type() != expectedType {
		t.Fatalf("Type() = %v, want %v", templateInclude.Type(), expectedType)
	}
}

func TestTemplateInclude_Equals(t *testing.T) {
	tests := []struct {
		name     string
		include1 *TemplateInclude
		include2 Include
		expected bool
	}{
		{
			name: "equal includes",
			include1: &TemplateInclude{
				Template: "my-template",
			},
			include2: &TemplateInclude{
				Template: "my-template",
			},
			expected: true,
		},
		{
			name: "different templates",
			include1: &TemplateInclude{
				Template: "my-template-1",
			},
			include2: &TemplateInclude{
				Template: "my-template-2",
			},
			expected: false,
		},
		{
			name: "different types",
			include1: &TemplateInclude{
				Template: "my-template",
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

func TestNewTemplateInclude(t *testing.T) {
	node := &yaml.Node{}
	template := "my-template"
	templateInclude := NewTemplateInclude(node, template)

	if templateInclude.Template != template {
		t.Fatalf("NewTemplateInclude().Template = %v, want %v", templateInclude.Template, template)
	}

	if templateInclude.Node() != node {
		t.Fatalf("NewTemplateInclude().Node() = %v, want %v", templateInclude.Node(), node)
	}
}
