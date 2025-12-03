package ci_yaml

import (
	"github.com/timo-reymann/gitlab-ci-verify/v2/pkg/location"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"strings"
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
		projectRoot + "/.gitlab-ci.yml",
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

func TestCreateVirtualCiYamlFile_EmptyGlobWarning(t *testing.T) {
	// Create a minimal test scenario
	projectRoot := "test_data/virtual-file"
	entryFileContent := []byte(`
include:
  - local: '.gitlab/ci/non-existent/*.yml'
`)

	entryFile, err := NewCiYamlFile(entryFileContent)
	if err != nil {
		t.Fatalf("Failed to create entry CiYamlFile: %v", err)
	}

	virtualFile, err := CreateVirtualCiYamlFile(projectRoot, projectRoot+"/.gitlab-ci.yml", entryFile)
	if err != nil {
		t.Fatalf("CreateVirtualCiYamlFile() error = %v", err)
	}

	// Check that we got a warning for the empty glob
	if len(virtualFile.Warnings) != 1 {
		t.Fatalf("Expected 1 warning, got %d", len(virtualFile.Warnings))
	}

	if virtualFile.Warnings[0].Code != 101 {
		t.Errorf("Expected code 101, got %d", virtualFile.Warnings[0].Code)
	}

	expectedMessage := "Include pattern '.gitlab/ci/non-existent/*.yml' did not match any files"
	if virtualFile.Warnings[0].Message != expectedMessage {
		t.Errorf("Expected '%s', got '%s'", expectedMessage, virtualFile.Warnings[0].Message)
	}
}

func TestCreateVirtualCiYamlFile_FileLoadError(t *testing.T) {
	// Create a test scenario with a non-existent single file
	projectRoot := "test_data/virtual-file"
	entryFileContent := []byte(`
include:
  - local: '.gitlab/ci/non-existent-file.yml'
`)

	entryFile, err := NewCiYamlFile(entryFileContent)
	if err != nil {
		t.Fatalf("Failed to create entry CiYamlFile: %v", err)
	}

	virtualFile, err := CreateVirtualCiYamlFile(projectRoot, projectRoot+"/.gitlab-ci.yml", entryFile)
	if err != nil {
		t.Fatalf("CreateVirtualCiYamlFile() error = %v", err)
	}

	// Check that we got an error warning for the missing file
	if len(virtualFile.Warnings) != 1 {
		t.Fatalf("Expected 1 warning, got %d", len(virtualFile.Warnings))
	}

	if virtualFile.Warnings[0].Code != 102 {
		t.Errorf("Expected code 102, got %d", virtualFile.Warnings[0].Code)
	}

	// Message should contain error details
	if !strings.Contains(virtualFile.Warnings[0].Message, "Include file") || !strings.Contains(virtualFile.Warnings[0].Message, "could not be loaded") {
		t.Errorf("Expected error message to contain file path and error details, got '%s'", virtualFile.Warnings[0].Message)
	}
}

func TestGetIgnoredCodes(t *testing.T) {
	tests := []struct {
		name        string
		virtualFile *VirtualCiYamlFile
		line        int
		expected    []string
	}{
		{
			name: "file level ignore",
			virtualFile: &VirtualCiYamlFile{
				Parts: []*VirtualCiYamlFilePart{
					{
						CiYamlFile: &CiYamlFile{
							ParsedYamlDoc: &yaml.Node{
								Content: []*yaml.Node{
									{
										Content: []*yaml.Node{
											{HeadComment: "gitlab-ci-verify ignore:CODE1"},
										},
									},
								},
							},
						},
						StartLine: 0,
						EndLine:   10,
					},
				},
			},
			line:     1,
			expected: []string{"CODE1"},
		},
		{
			name: "line level ignore",
			virtualFile: &VirtualCiYamlFile{
				Parts: []*VirtualCiYamlFilePart{
					{
						CiYamlFile: &CiYamlFile{
							lineNumberMapping: yaml.LineNumberMapping{
								2: []*yaml.Node{
									{LineComment: "gitlab-ci-verify ignore:CODE2"},
								},
							},
							ParsedYamlDoc: &yaml.Node{},
						},
						StartLine: 1,
						EndLine:   10,
					},
				},
			},
			line:     2,
			expected: []string{"CODE2"},
		},
		{
			name: "file and line level ignore",
			virtualFile: &VirtualCiYamlFile{
				Parts: []*VirtualCiYamlFilePart{
					{
						CiYamlFile: &CiYamlFile{
							ParsedYamlDoc: &yaml.Node{
								Content: []*yaml.Node{
									{
										Content: []*yaml.Node{
											{HeadComment: "gitlab-ci-verify ignore:CODE1"},
										},
									},
								},
							},
							lineNumberMapping: yaml.LineNumberMapping{
								2: []*yaml.Node{
									{LineComment: "gitlab-ci-verify ignore:CODE2"},
								},
							},
						},
						StartLine: 1,
						EndLine:   10,
					},
				},
			},
			line:     2,
			expected: []string{"CODE1", "CODE2"},
		},
		{
			name: "no ignore",
			virtualFile: &VirtualCiYamlFile{
				Parts: []*VirtualCiYamlFilePart{
					{
						CiYamlFile: &CiYamlFile{
							ParsedYamlDoc: &yaml.Node{},
						},
						StartLine: 1,
						EndLine:   10,
					},
				},
			},
			line:     1,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.virtualFile.GetIgnoredCodes(tt.line)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetIgnoredCodes(%d) = %v, want %v", tt.line, result, tt.expected)
			}
		})
	}
}
