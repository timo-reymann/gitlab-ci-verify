package ci_yaml

import (
	"os"
	"testing"
)

func TestCreateVirtualCiYamlFile(t *testing.T) {
	projectRoot := "test_data/virtual-file"
	entryFilePath := projectRoot + "/.gitlab-ci.yml"

	entryFileContent, err := os.ReadFile(entryFilePath)
	if err != nil {
		t.Fatalf("Failed to read entry file: %v", err)
	}

	entryFile, err := NewCiYamlFile(entryFileContent)
	if err != nil {
		t.Fatalf("Failed to create entry CiYamlFile: %v", err)
	}

	virtualFile, err := CreateVirtualCiYamlFile(projectRoot, entryFilePath, entryFile)
	if err != nil {
		t.Fatalf("CreateVirtualCiYamlFile() error = %v", err)
	}

	expectedParts := []string{
		projectRoot + "/.gitlab/ci/templates/.mod_download.gitlab-ci.yml",
		projectRoot + "/.gitlab/ci/pipelines/release.gitlab-ci.yml",
		projectRoot + "/.gitlab/ci/pipelines/merge_request.gitlab-ci.yml",
		projectRoot + "/.gitlab/ci/pipelines/main.gitlab-ci.yml",
	}

	if len(virtualFile.Parts) != len(expectedParts) {
		t.Fatalf("Expected %d parts, got %d", len(expectedParts), len(virtualFile.Parts))
	}

	for i, part := range expectedParts {
		if part != virtualFile.Parts[i].Path {
			t.Errorf("Expected part path to be %s, got %s", part, virtualFile.Parts[i].Path)
		}
	}

	for i, part := range virtualFile.Parts {
		if part.Path != expectedParts[i] {
			t.Errorf("Expected part path to be %s, got %s", expectedParts[i], part.Path)
		}
	}
}
