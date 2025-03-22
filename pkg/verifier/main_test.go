package verifier

import (
	"bytes"
	"github.com/timo-reymann/gitlab-ci-verify/internal/cli"
	ci_yaml "github.com/timo-reymann/gitlab-ci-verify/internal/gitlab/ci-yaml"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/checks"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/formatter"
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

type DummyCheck struct {
	name     string
	findings []checks.CheckFinding
}

func (dc *DummyCheck) Name() string {
	return dc.name
}

func (dc *DummyCheck) Run(input *checks.CheckInput) ([]checks.CheckFinding, error) {
	return dc.findings, nil
}

func TestRunChecks(t *testing.T) {
	tests := []struct {
		name         string
		checkInput   *checks.CheckInput
		failSeverity int
		expectedExit bool
		checks       []checks.Check
	}{
		{
			name: "no findings",
			checkInput: &checks.CheckInput{
				VirtualCiYaml: &ci_yaml.VirtualCiYamlFile{},
				Configuration: &cli.Configuration{},
			},
			failSeverity: checks.SeverityError,
			expectedExit: false,
			checks: []checks.Check{
				&DummyCheck{name: "DummyCheck1", findings: []checks.CheckFinding{}},
			},
		},
		{
			name: "findings below fail severity",
			checkInput: &checks.CheckInput{
				VirtualCiYaml: &ci_yaml.VirtualCiYamlFile{},
				Configuration: &cli.Configuration{},
			},
			failSeverity: checks.SeverityError,
			expectedExit: false,
			checks: []checks.Check{
				&DummyCheck{name: "DummyCheck2", findings: []checks.CheckFinding{
					{Severity: checks.SeverityWarning},
				}},
			},
		},
		{
			name: "findings at fail severity",
			checkInput: &checks.CheckInput{
				VirtualCiYaml: &ci_yaml.VirtualCiYamlFile{},
				Configuration: &cli.Configuration{},
			},
			failSeverity: checks.SeverityWarning,
			expectedExit: true,
			checks: []checks.Check{
				&DummyCheck{name: "DummyCheck3", findings: []checks.CheckFinding{
					{Severity: checks.SeverityWarning},
				}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gcv := &GitlabCIVerifier{
				configuration: &cli.Configuration{
					FailSeverity: "error",
				},
				formatter: &formatter.TextFindingsFormatter{},
			}

			writer := &bytes.Buffer{}
			err := gcv.SetupFormatter(writer, "text")
			if err != nil {
				t.Fatal(err)
			}

			errorHandler := func(err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			isExit := gcv.RunChecks(tt.checkInput, tt.checks, tt.failSeverity, errorHandler)
			if isExit != tt.expectedExit {
				t.Errorf("RunChecks() = %v, want %v", isExit, tt.expectedExit)
			}
		})
	}
}
