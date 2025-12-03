package includes

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
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

func TestLocalInclude_IsGlobPattern(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "simple path",
			path:     "test.yml",
			expected: false,
		},
		{
			name:     "asterisk wildcard",
			path:     "*.yml",
			expected: true,
		},
		{
			name:     "directory with asterisk",
			path:     ".gitlab/ci/*.yml",
			expected: true,
		},
		{
			name:     "question mark wildcard",
			path:     "test?.yml",
			expected: true,
		},
		{
			name:     "bracket wildcard",
			path:     "test[0-9].yml",
			expected: true,
		},
		{
			name:     "path without wildcards",
			path:     "/absolute/path/test.yml",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			localInclude := &LocalInclude{Path: tt.path}
			result := localInclude.IsGlobPattern()
			if result != tt.expected {
				t.Fatalf("IsGlobPattern() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLocalInclude_ResolvePaths(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	ciDir := filepath.Join(tempDir, ".gitlab", "ci")
	err := os.MkdirAll(ciDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create test files
	testFiles := []string{"a.yml", "b.yml", "c.yml"}
	for _, file := range testFiles {
		filePath := filepath.Join(ciDir, file)
		err := os.WriteFile(filePath, []byte("test: content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	tests := []struct {
		name          string
		path          string
		srcFile       string
		expectedCount int
		expectError   bool
	}{
		{
			name:          "non-glob path",
			path:          ".gitlab/ci/a.yml",
			srcFile:       ".gitlab-ci.yml",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "glob pattern matching multiple files",
			path:          ".gitlab/ci/*.yml",
			srcFile:       ".gitlab-ci.yml",
			expectedCount: 3,
			expectError:   false,
		},
		{
			name:          "glob pattern matching no files",
			path:          ".gitlab/ci/*.json",
			srcFile:       ".gitlab-ci.yml",
			expectedCount: 0,
			expectError:   false,
		},
		{
			name:          "glob pattern with question mark",
			path:          ".gitlab/ci/?.yml",
			srcFile:       ".gitlab-ci.yml",
			expectedCount: 3,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			localInclude := &LocalInclude{Path: tt.path}
			paths, err := localInclude.ResolvePaths(tempDir, tt.srcFile)

			if tt.expectError && err == nil {
				t.Fatalf("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(paths) != tt.expectedCount {
				t.Fatalf("ResolvePaths() returned %d paths, want %d", len(paths), tt.expectedCount)
			}

			// Verify paths are sorted
			for i := 1; i < len(paths); i++ {
				if paths[i-1] >= paths[i] {
					t.Fatalf("ResolvePaths() returned unsorted paths: %v", paths)
				}
			}
		})
	}
}
