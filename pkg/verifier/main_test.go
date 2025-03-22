package verifier

import (
	"bytes"
	"github.com/timo-reymann/gitlab-ci-verify/internal/cli"
	"os"
	"testing"
)

func TestShouldCheckAgainstLintAPI(t *testing.T) {
	tests := []struct {
		name          string
		configuration *cli.Configuration
		expected      bool
		ci            bool
	}{
		{
			name: "CI environment with no lint API call",
			configuration: &cli.Configuration{
				NoLintAPICallInCi: true,
			},
			ci:       true,
			expected: false,
		},
		{
			name: "CI environment with lint API call",
			configuration: &cli.Configuration{
				NoLintAPICallInCi: false,
			},
			ci:       true,
			expected: true,
		},
		{
			name:          "Non-CI environment",
			configuration: &cli.Configuration{},
			ci:            false,
			expected:      true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if tt.ci {
				os.Setenv("CI", "true")
			} else {
				os.Unsetenv("CI")
			}
			gcv := &GitlabCIVerifier{
				configuration: tt.configuration,
			}
			result := gcv.shouldCheckAgainstLintAPI()
			if result != tt.expected {
				t.Errorf("shouldCheckAgainstLintAPI() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSetupFormatter(t *testing.T) {
	tests := []struct {
		name          string
		formatterName string
		expectedError bool
	}{
		{
			name:          "valid formatter",
			formatterName: "text",
			expectedError: false,
		},
		{
			name:          "invalid formatter",
			formatterName: "invalid",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gcv := &GitlabCIVerifier{
				configuration: &cli.Configuration{
					FailSeverity: "error",
				},
			}

			writer := &bytes.Buffer{}
			err := gcv.SetupFormatter(writer, tt.formatterName)
			if (err != nil) != tt.expectedError {
				t.Errorf("SetupFormatter() error = %v, wantErr %v", err, tt.expectedError)
			}
		})
	}
}
