package checks

import (
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/gitlab/ci-yaml"
	"testing"
)

func TestLocalIncludeGlobCheck_Run_NoWarnings(t *testing.T) {
	check := LocalIncludeGlobCheck{}
	input := &CheckInput{
		VirtualCiYaml: &ci_yaml.VirtualCiYamlFile{
			EntryFilePath: "test.yml",
			Warnings:      []ci_yaml.VirtualFileWarning{},
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

func TestLocalIncludeGlobCheck_Run_WithWarnings(t *testing.T) {
	check := LocalIncludeGlobCheck{}
	input := &CheckInput{
		VirtualCiYaml: &ci_yaml.VirtualCiYamlFile{
			EntryFilePath: "test.yml",
			Warnings: []ci_yaml.VirtualFileWarning{
				{
					Message:     "Glob pattern did not match any files",
					IncludePath: ".gitlab/ci/*.yml",
				},
				{
					Message:     "Glob pattern did not match any files",
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
