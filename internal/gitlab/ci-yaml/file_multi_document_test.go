package ci_yaml

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewCiYamlFile_MultiDocument(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		expectHeader bool
		expectedKeys []string
	}{
		{
			name:         "single document",
			content:      "build:\n  script: echo\n",
			expectHeader: false,
			expectedKeys: []string{"build"},
		},
		{
			name:         "header and configuration",
			content:      "spec:\n  inputs:\n    stage:\n---\nbuild:\n  stage: $[[ inputs.stage ]]\n  script: echo\n",
			expectHeader: true,
			expectedKeys: []string{"build"},
		},
		{
			name:         "leading document marker",
			content:      "---\nbuild:\n  script: echo\n",
			expectHeader: false,
			expectedKeys: []string{"build"},
		},
		{
			name:         "leading document marker with header",
			content:      "---\nspec:\n  inputs: {}\n---\nbuild:\n  script: echo\n",
			expectHeader: true,
			expectedKeys: []string{"build"},
		},
		{
			name:         "two documents without spec use the first one",
			content:      "build:\n  script: echo\n---\nother:\n  script: echo\n",
			expectHeader: false,
			expectedKeys: []string{"build"},
		},
		{
			name:         "spec only",
			content:      "spec:\n  inputs: {}\n",
			expectHeader: false,
			expectedKeys: []string{"spec"},
		},
		{
			name:         "invalid content after first document is ignored",
			content:      "{}\ninclude: []\n",
			expectHeader: false,
			expectedKeys: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewCiYamlFile([]byte(tt.content))
			if err != nil {
				t.Fatalf("NewCiYamlFile() error = %v", err)
			}

			if f.HasHeader() != tt.expectHeader {
				t.Errorf("HasHeader() = %v, want %v", f.HasHeader(), tt.expectHeader)
			}

			if len(f.ParsedYamlMap) != len(tt.expectedKeys) {
				t.Fatalf("ParsedYamlMap = %v, want keys %v", f.ParsedYamlMap, tt.expectedKeys)
			}
			for _, key := range tt.expectedKeys {
				if _, ok := f.ParsedYamlMap[key]; !ok {
					t.Errorf("ParsedYamlMap = %v, missing key %s", f.ParsedYamlMap, key)
				}
			}

			root := f.ParsedYamlDoc.Content[0]
			if len(tt.expectedKeys) > 0 && root.Content[0].Value != tt.expectedKeys[0] {
				t.Errorf("ParsedYamlDoc first key = %s, want %s", root.Content[0].Value, tt.expectedKeys[0])
			}
		})
	}
}

func TestNewCiYamlFile_MultiDocumentKeepsLineNumbers(t *testing.T) {
	f, err := NewCiYamlFile([]byte("spec:\n  inputs: {}\n---\nbuild:\n  script:\n    - echo $FOO # gitlab-ci-verify ignore:SC-2086\n"))
	if err != nil {
		t.Fatal(err)
	}

	scriptItem := f.ParsedYamlDoc.Content[0].Content[1].Content[1].Content[0]
	if scriptItem.Line != 6 {
		t.Errorf("script line = %d, want 6", scriptItem.Line)
	}

	ignores := f.GetLineLevelIgnores(6)
	if len(ignores) != 1 || ignores[0].Code != "SC-2086" {
		t.Errorf("GetLineLevelIgnores(6) = %v, want SC-2086", ignores)
	}
}

func TestGetFileLevelIgnores_WithHeader(t *testing.T) {
	f, err := NewCiYamlFile([]byte("# gitlab-ci-verify ignore:SC-2086\nspec:\n  inputs: {}\n---\nbuild:\n  script: echo $FOO\n"))
	if err != nil {
		t.Fatal(err)
	}

	ignores := f.GetFileLevelIgnores()
	if len(ignores) != 1 || ignores[0].Code != "SC-2086" {
		t.Errorf("GetFileLevelIgnores() = %v, want SC-2086", ignores)
	}
}

func TestCiYamlFile_ConfigurationContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "single document is unchanged",
			content:  "build:\n  script: echo\n",
			expected: "build:\n  script: echo\n",
		},
		{
			name:     "header is blanked",
			content:  "spec:\n  inputs: {}\n---\n# comment\nbuild:\n  script: echo\n",
			expected: "\n\n\n# comment\nbuild:\n  script: echo\n",
		},
		{
			name:     "leading document marker is blanked",
			content:  "---\nbuild:\n  script: echo\n",
			expected: "\nbuild:\n  script: echo\n",
		},
		{
			name:     "following documents are blanked",
			content:  "build:\n  script: echo\n---\nother: {}\n",
			expected: "build:\n  script: echo\n\n\n",
		},
		{
			name:     "document marker inside block scalar is kept",
			content:  "build:\n  script: |\n    echo\n    ---\n",
			expected: "build:\n  script: |\n    echo\n    ---\n",
		},
		{
			name:     "windows line endings",
			content:  "spec:\r\n  inputs: {}\r\n---\r\nbuild:\r\n  script: echo\r\n",
			expected: "\n\n\nbuild:\r\n  script: echo\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewCiYamlFile([]byte(tt.content))
			if err != nil {
				t.Fatal(err)
			}

			actual := f.ConfigurationContent()
			if string(actual) != tt.expected {
				t.Errorf("ConfigurationContent() = %q, want %q", actual, tt.expected)
			}

			if bytes.Count(actual, []byte("\n")) != bytes.Count(f.FileContent, []byte("\n")) {
				t.Errorf("ConfigurationContent() changed the line count")
			}
		})
	}
}

func TestCreateVirtualCiYamlFile_IncludesWithMultipleDocuments(t *testing.T) {
	projectRoot := t.TempDir()
	writeFile := func(name, content string) {
		if err := os.WriteFile(filepath.Join(projectRoot, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	writeFile("component.yml", "spec:\n  inputs:\n    stage:\n---\ncomponent-job:\n  stage: $[[ inputs.stage ]]\n  script: echo\n")
	writeFile("plain.yml", "---\nplain-job:\n  script: echo\n")

	entryFilePath := filepath.Join(projectRoot, ".gitlab-ci.yml")
	entryContent := "include:\n  - local: component.yml\n    inputs:\n      stage: test\n  - local: plain.yml\n\nentry-job:\n  script: echo\n"
	writeFile(".gitlab-ci.yml", entryContent)

	entryFile, err := NewCiYamlFile([]byte(entryContent))
	if err != nil {
		t.Fatal(err)
	}

	virtualFile, err := CreateVirtualCiYamlFile(projectRoot, entryFilePath, entryFile)
	if err != nil {
		t.Fatal(err)
	}

	for _, job := range []string{"entry-job", "component-job", "plain-job"} {
		if _, ok := virtualFile.Combined.ParsedYamlMap[job]; !ok {
			t.Errorf("Combined is missing job %s:\n%s", job, virtualFile.Combined.FileContent)
		}
	}

	lintYaml := string(virtualFile.LintYaml)
	if strings.Contains(lintYaml, "component-job") || strings.Contains(lintYaml, "spec:") {
		t.Errorf("LintYaml must not inline files with a header:\n%s", lintYaml)
	}
	if !strings.Contains(lintYaml, "plain-job") || !strings.Contains(lintYaml, "local: component.yml") {
		t.Errorf("LintYaml is missing plain include or the include block:\n%s", lintYaml)
	}

	for idx, part := range virtualFile.Parts {
		if idx > 0 && virtualFile.Parts[idx-1].EndLine >= part.StartLine {
			t.Errorf("Part %s overlaps with the previous part", part.Path)
		}

		_, loc := virtualFile.Resolve(part.StartLine)
		if loc == nil || loc.File != part.Path || loc.Line != 0 {
			t.Errorf("Resolve(%d) = %v, want %s:0", part.StartLine, loc, part.Path)
		}
	}

	componentPart := virtualFile.Parts[1]
	componentJob := virtualFile.Combined.ParsedYamlDoc.Content[0]
	var componentJobLine int
	for idx := 0; idx < len(componentJob.Content); idx += 2 {
		if componentJob.Content[idx].Value == "component-job" {
			componentJobLine = componentJob.Content[idx].Line
		}
	}
	_, loc := virtualFile.Resolve(componentJobLine)
	if loc == nil || loc.File != componentPart.Path || loc.Line != 4 {
		t.Errorf("Resolve(%d) = %v, want %s:4", componentJobLine, loc, componentPart.Path)
	}
}
