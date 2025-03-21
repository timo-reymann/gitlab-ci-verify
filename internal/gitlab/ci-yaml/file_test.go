package ci_yaml

import (
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
