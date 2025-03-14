package ci_yaml

import (
	"github.com/timo-reymann/gitlab-ci-verify/pkg/location"
	"os"
	"reflect"
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

func TestVirtualCiYamlFile_Resolve(t *testing.T) {
	tests := []struct {
		name        string
		parts       []*VirtualCiYamlFilePart
		line        int
		expected    *VirtualCiYamlFilePart
		expectedLoc *location.Location
	}{
		{
			name: "line in first part",
			parts: []*VirtualCiYamlFilePart{
				{StartLine: 1, EndLine: 10, Path: "part1.yml"},
				{StartLine: 11, EndLine: 20, Path: "part2.yml"},
			},
			line:        5,
			expected:    &VirtualCiYamlFilePart{StartLine: 1, EndLine: 10, Path: "part1.yml"},
			expectedLoc: location.NewLocation("part1.yml", 4),
		},
		{
			name: "line in second part",
			parts: []*VirtualCiYamlFilePart{
				{StartLine: 1, EndLine: 10, Path: "part1.yml"},
				{StartLine: 11, EndLine: 20, Path: "part2.yml"},
			},
			line:        15,
			expected:    &VirtualCiYamlFilePart{StartLine: 11, EndLine: 20, Path: "part2.yml"},
			expectedLoc: location.NewLocation("part2.yml", 4),
		},
		{
			name: "line before first part",
			parts: []*VirtualCiYamlFilePart{
				{StartLine: 1, EndLine: 10, Path: "part1.yml"},
				{StartLine: 11, EndLine: 20, Path: "part2.yml"},
			},
			line:        0,
			expected:    nil,
			expectedLoc: nil,
		},
		{
			name: "line after last part",
			parts: []*VirtualCiYamlFilePart{
				{StartLine: 1, EndLine: 10, Path: "part1.yml"},
				{StartLine: 11, EndLine: 20, Path: "part2.yml"},
			},
			line:        21,
			expected:    nil,
			expectedLoc: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VirtualCiYamlFile{
				Parts: tt.parts,
			}
			result, loc := v.Resolve(tt.line)
			if !reflect.DeepEqual(tt.expected, result) || !reflect.DeepEqual(tt.expectedLoc, loc) {
				t.Errorf("Resolve(%d) = %v, %v; want %v, %v", tt.line, result, loc, tt.expected, tt.expectedLoc)
			}
		})
	}
}
