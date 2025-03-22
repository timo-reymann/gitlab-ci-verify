package checks

import (
	ciyaml "github.com/timo-reymann/gitlab-ci-verify/internal/gitlab/ci-yaml"
	"testing"
)

func TestCheckInput_HasLintAPIResult(t *testing.T) {
	tests := []struct {
		name     string
		input    *CheckInput
		expected bool
	}{
		{
			name: "with LintAPIResult",
			input: &CheckInput{
				LintAPIResult: &ciyaml.VerificationResultWithRemoteInfo{},
			},
			expected: true,
		},
		{
			name: "without LintAPIResult",
			input: &CheckInput{
				LintAPIResult: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.HasLintAPIResult()
			if result != tt.expected {
				t.Errorf("HasLintAPIResult() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCheckInput_CanProvideMergedYaml(t *testing.T) {
	tests := []struct {
		name     string
		input    *CheckInput
		expected bool
	}{
		{
			name: "without LintAPIResult",
			input: &CheckInput{
				LintAPIResult: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.CanProvideMergedYaml()
			if result != tt.expected {
				t.Errorf("CanProvideMergedYaml() = %v, want %v", result, tt.expected)
			}
		})
	}
}
