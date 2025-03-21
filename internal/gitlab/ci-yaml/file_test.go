package ci_yaml

import (
	"github.com/timo-reymann/gitlab-ci-verify/pkg/filtering"
	"os"
	"testing"
)

func TestNewCiYamlFile(t *testing.T) {
	// Read a sample YAML file
	content, err := os.ReadFile("test_data/gitlab-ci/validSingleJobSingleLine.yaml")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	// Create a new CiYamlFile
	ciYamlFile, err := NewCiYamlFile(content)
	if err != nil {
		t.Fatalf("NewCiYamlFile() error = %v", err)
	}

	// Check if the file content matches
	if string(ciYamlFile.FileContent) != string(content) {
		t.Fatalf("FileContent = %v, want %v", string(ciYamlFile.FileContent), string(content))
	}

	// Check if the parsed YAML map is not nil
	if ciYamlFile.ParsedYamlMap == nil {
		t.Fatal("ParsedYamlMap is nil")
	}

	// Check if the parsed YAML document is not nil
	if ciYamlFile.ParsedYamlDoc == nil {
		t.Fatal("ParsedYamlDoc is nil")
	}
}

func TestGetFileLevelIgnores(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent string
		expected    []filtering.IgnoreComment
	}{
		{
			name: "single ignore code",
			yamlContent: `# gitlab-ci-verify ignore:CODE1
stages:
  - build
  - test
`,
			expected: []filtering.IgnoreComment{
				{Comment: "gitlab-ci-verify ignore:CODE1", Code: "CODE1"},
			},
		},
		{
			name: "multiple ignore codes",
			yamlContent: `# gitlab-ci-verify ignore:CODE1 ignore:CODE2
stages:
  - build
  - test
`,
			expected: []filtering.IgnoreComment{
				{Comment: "gitlab-ci-verify ignore:CODE1 ignore:CODE2", Code: "CODE1"},
				{Comment: "gitlab-ci-verify ignore:CODE1 ignore:CODE2", Code: "CODE2"},
			},
		},
		{
			name: "multiple line ignore codes",
			yamlContent: `# gitlab-ci-verify ignore:CODE1 ignore:CODE2
# gitlab-ci-verify ignore:CODE3
stages:
  - build
  - test
`,
			expected: []filtering.IgnoreComment{
				{Comment: "gitlab-ci-verify ignore:CODE1 ignore:CODE2", Code: "CODE1"},
				{Comment: "gitlab-ci-verify ignore:CODE1 ignore:CODE2", Code: "CODE2"},
				{Comment: "gitlab-ci-verify ignore:CODE3", Code: "CODE3"},
			},
		},
		{
			name: "no ignore code",
			yamlContent: `# some other comment
stages:
  - build
  - test
`,
			expected: []filtering.IgnoreComment{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ciYamlFile, err := NewCiYamlFile([]byte(tt.yamlContent))
			if err != nil {
				t.Fatalf("NewCiYamlFile() error = %v", err)
			}

			result := ciYamlFile.GetFileLevelIgnores()
			if len(result) != len(tt.expected) {
				t.Errorf("GetFileLevelIgnores() = %v, want %v", result, tt.expected)
			}
			for i := range result {
				if result[i].Comment != tt.expected[i].Comment || result[i].Code != tt.expected[i].Code {
					t.Errorf("GetFileLevelIgnores() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestGetLineLevelIgnores(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent string
		line        int
		expected    []filtering.IgnoreComment
	}{
		{
			name: "single line with ignore comment",
			yamlContent: `
stages: # gitlab-ci-verify ignore:CODE1
  - build
  - test
`,
			line: 2,
			expected: []filtering.IgnoreComment{
				{Comment: "gitlab-ci-verify ignore:CODE1", Code: "CODE1"},
			},
		},
		{
			name: "multiple lines with ignore comments",
			yamlContent: `
stages: # gitlab-ci-verify ignore:CODE1
  - build # gitlab-ci-verify ignore:CODE2
  - test
`,
			line: 3,
			expected: []filtering.IgnoreComment{
				{Comment: "gitlab-ci-verify ignore:CODE2", Code: "CODE2"},
			},
		},
		{
			name: "line without ignore comment",
			yamlContent: `
stages:
  - build
  - test
`,
			line:     1,
			expected: []filtering.IgnoreComment{},
		},
		{
			name: "multiple nodes with ignore comments",
			yamlContent: `
stages: # gitlab-ci-verify ignore:CODE1
  - build # gitlab-ci-verify ignore:CODE2
  - test # gitlab-ci-verify ignore:CODE3
`,
			line: 3,
			expected: []filtering.IgnoreComment{
				{Comment: "gitlab-ci-verify ignore:CODE2", Code: "CODE2"},
			},
		},
		{
			name: "multiple nodes with ignore comments",
			yamlContent: `
# gitlab-ci-verify ignore:CODE1
stages: 
  - build
  - test
`,
			line: 3,
			expected: []filtering.IgnoreComment{
				{Comment: "gitlab-ci-verify ignore:CODE1", Code: "CODE1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ciYamlFile, err := NewCiYamlFile([]byte(tt.yamlContent))
			if err != nil {
				t.Fatalf("NewCiYamlFile() error = %v", err)
			}

			result := ciYamlFile.GetLineLevelIgnores(tt.line)
			if len(result) != len(tt.expected) {
				t.Errorf("GetLineLevelIgnores() = %v, want %v", result, tt.expected)
			}
			for i := range result {
				if result[i].Comment != tt.expected[i].Comment || result[i].Code != tt.expected[i].Code {
					t.Errorf("GetLineLevelIgnores() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}
