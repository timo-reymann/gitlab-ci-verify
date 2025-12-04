package checks

import (
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/gitlab/ci-yaml"
	"testing"
)

func TestLocalIncludeGlobCheck_Run_NoResolveFindings(t *testing.T) {
	check := LocalIncludeGlobCheck{}
	input := &CheckInput{
		VirtualCiYaml: &ci_yaml.VirtualCiYamlFile{
			EntryFilePath: "test.yml",
			ResolveFindings:      []ci_yaml.VirtualFileResolveFinding{},
		},
	}

	findings, err := check.Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(findings) != 0 {
		t.Fatalf("expected 0 findings, got %d", len(findings))
	}
}

func TestLocalIncludeGlobCheck_Run_WithResolveFindings(t *testing.T) {
	check := LocalIncludeGlobCheck{}
	input := &CheckInput{
		VirtualCiYaml: &ci_yaml.VirtualCiYamlFile{
			EntryFilePath: "test.yml",
			ResolveFindings: []ci_yaml.VirtualFileResolveFinding{
				{
					Code:        101,
					Severity:    1,
					Message:     "Include pattern '.gitlab/ci/*.yml' did not match any files",
					IncludePath: ".gitlab/ci/*.yml",
				},
				{
					Code:        101,
					Severity:    1,
					Message:     "Include pattern 'includes/**/*.yaml' did not match any files",
					IncludePath: "includes/**/*.yaml",
				},
			},
		},
	}

	findings, err := check.Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(findings))
	}

	// Check first finding
	if findings[0].Severity != SeverityWarning {
		t.Errorf("expected severity %d, got %d", SeverityWarning, findings[0].Severity)
	}
	if findings[0].Code != "INC-101" {
		t.Errorf("expected code INC-101, got %s", findings[0].Code)
	}
	expectedMsg := "Include pattern '.gitlab/ci/*.yml' did not match any files"
	if findings[0].Message != expectedMsg {
		t.Errorf("expected message '%s', got '%s'", expectedMsg, findings[0].Message)
	}

	// Check second finding
	if findings[1].Severity != SeverityWarning {
		t.Errorf("expected severity %d, got %d", SeverityWarning, findings[1].Severity)
	}
	expectedMsg2 := "Include pattern 'includes/**/*.yaml' did not match any files"
	if findings[1].Message != expectedMsg2 {
		t.Errorf("expected message '%s', got '%s'", expectedMsg2, findings[1].Message)
	}
}

func TestLocalIncludeGlobCheck_Run_WithFileLoadErrors(t *testing.T) {
	check := LocalIncludeGlobCheck{}
	input := &CheckInput{
		VirtualCiYaml: &ci_yaml.VirtualCiYamlFile{
			EntryFilePath: "test.yml",
			ResolveFindings: []ci_yaml.VirtualFileResolveFinding{
				{
					Code:        102,
					Severity:    0,
					Message:     "Include file 'missing.yml' could not be loaded: open test.yml: no such file or directory",
					IncludePath: "missing.yml",
				},
			},
		},
	}

	findings, err := check.Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}

	// Check finding is ERROR severity
	if findings[0].Severity != SeverityError {
		t.Errorf("expected severity %d (Error), got %d", SeverityError, findings[0].Severity)
	}
	if findings[0].Code != "INC-102" {
		t.Errorf("expected code INC-102, got %s", findings[0].Code)
	}
	expectedMsg := "Include file 'missing.yml' could not be loaded: open test.yml: no such file or directory"
	if findings[0].Message != expectedMsg {
		t.Errorf("expected message '%s', got '%s'", expectedMsg, findings[0].Message)
	}
}

func TestLocalIncludeGlobCheck_Run_WithMixedResolveFindingsAndErrors(t *testing.T) {
	check := LocalIncludeGlobCheck{}
	input := &CheckInput{
		VirtualCiYaml: &ci_yaml.VirtualCiYamlFile{
			EntryFilePath: "test.yml",
			ResolveFindings: []ci_yaml.VirtualFileResolveFinding{
				{
					Code:        101,
					Severity:    1,
					Message:     "Include pattern '.gitlab/ci/*.yml' did not match any files",
					IncludePath: ".gitlab/ci/*.yml",
				},
				{
					Code:        102,
					Severity:    0,
					Message:     "Include file 'protected.yml' could not be loaded: permission denied",
					IncludePath: "protected.yml",
				},
			},
		},
	}

	findings, err := check.Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(findings))
	}

	// First should be warning
	if findings[0].Severity != SeverityWarning {
		t.Errorf("expected first finding to be Warning severity, got %d", findings[0].Severity)
	}

	// Second should be error
	if findings[1].Severity != SeverityError {
		t.Errorf("expected second finding to be Error severity, got %d", findings[1].Severity)
	}
}
