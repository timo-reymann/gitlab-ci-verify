package includes

import (
	"gopkg.in/yaml.v3"
	"slices"
	"testing"
)

func TestProjectInclude_Type(t *testing.T) {
	projectInclude := &ProjectInclude{}
	expectedType := "project"
	if projectInclude.Type() != expectedType {
		t.Fatalf("Type() = %v, want %v", projectInclude.Type(), expectedType)
	}
}

func TestProjectInclude_Equals(t *testing.T) {
	tests := []struct {
		name     string
		include1 *ProjectInclude
		include2 Include
		expected bool
	}{
		{
			name: "equal includes",
			include1: &ProjectInclude{
				Project: "my-group/my-project",
				Files:   []string{"test.yml"},
			},
			include2: &ProjectInclude{
				Project: "my-group/my-project",
				Files:   []string{"test.yml"},
			},
			expected: true,
		},
		{
			name: "different projects",
			include1: &ProjectInclude{
				Project: "my-group/my-project-1",
				Files:   []string{"test.yml"},
			},
			include2: &ProjectInclude{
				Project: "my-group/my-project-2",
				Files:   []string{"test.yml"},
			},
			expected: false,
		},
		{
			name: "different files",
			include1: &ProjectInclude{
				Project: "my-group/my-project",
				Files:   []string{"test1.yml"},
			},
			include2: &ProjectInclude{
				Project: "my-group/my-project",
				Files:   []string{"test2.yml"},
			},
			expected: false,
		},
		{
			name: "different types",
			include1: &ProjectInclude{
				Project: "my-group/my-project",
				Files:   []string{"test.yml"},
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

func TestNewProjectInclude(t *testing.T) {
	node := &yaml.Node{}
	project := "my-group/my-project"
	files := []string{"test.yml"}
	projectInclude := NewProjectInclude(node, project, files)

	if projectInclude.Project != project {
		t.Fatalf("NewProjectInclude().Project = %v, want %v", projectInclude.Project, project)
	}

	if !slices.Equal(projectInclude.Files, files) {
		t.Fatalf("NewProjectInclude().Files = %v, want %v", projectInclude.Files, files)
	}

	if projectInclude.Node() != node {
		t.Fatalf("NewProjectInclude().Node() = %v, want %v", projectInclude.Node(), node)
	}
}
