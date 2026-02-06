package verifier

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/cli"
	ci_yaml "github.com/timo-reymann/gitlab-ci-verify/v2/internal/gitlab/ci-yaml"
	"github.com/timo-reymann/gitlab-ci-verify/v2/pkg/checks"
	"github.com/timo-reymann/gitlab-ci-verify/v2/pkg/formatter"
)

func TestShouldCheckAgainstLintAPI(t *testing.T) {
	tests := []struct {
		name          string
		configuration *cli.Configuration
		expected      bool
		ci            bool
	}{
		{
			name: "Online CI environment with no lint API call",
			configuration: &cli.Configuration{
				NoLintAPICallInCi: true,
			},
			ci:       true,
			expected: false,
		},
		{
			name: "Online CI environment with lint API call",
			configuration: &cli.Configuration{
				NoLintAPICallInCi: false,
			},
			ci:       true,
			expected: true,
		},
		{
			name:          "Online Non-CI environment",
			configuration: &cli.Configuration{},
			ci:            false,
			expected:      true,
		},
		{
			name: "Offline CI environment with no lint API call",
			configuration: &cli.Configuration{
				NoLintAPICallInCi: true,
				Offline:           true,
			},
			ci:       true,
			expected: false,
		},
		{
			name: "Offline CI environment with lint API call",
			configuration: &cli.Configuration{
				NoLintAPICallInCi: false,
				Offline:           true,
			},
			ci:       true,
			expected: false,
		},
		{
			name: "Offline Non-CI environment",
			configuration: &cli.Configuration{
				Offline: true,
			},
			ci:       false,
			expected: false,
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

			gcv := NewGitlabCIVerifier(&cli.Configuration{
				FailSeverity: "error",
			}, "/tmp")

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
		name                              string
		checkInput                        *checks.CheckInput
		failSeverity                      int
		expectedExit                      bool
		expectedReportLineCountWithHeader int
		checks                            []checks.Check
	}{
		{
			name: "no findings",
			checkInput: &checks.CheckInput{
				VirtualCiYaml: &ci_yaml.VirtualCiYamlFile{},
				Configuration: &cli.Configuration{},
			},
			failSeverity:                      checks.SeverityError,
			expectedExit:                      false,
			expectedReportLineCountWithHeader: 1,
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
			failSeverity:                      checks.SeverityError,
			expectedExit:                      false,
			expectedReportLineCountWithHeader: 2,
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
			failSeverity:                      checks.SeverityWarning,
			expectedExit:                      true,
			expectedReportLineCountWithHeader: 2,
			checks: []checks.Check{
				&DummyCheck{name: "DummyCheck3", findings: []checks.CheckFinding{
					{Severity: checks.SeverityWarning},
				}},
			},
		},
		{
			name: "multiple findings at the same line for the same reason",
			checkInput: &checks.CheckInput{
				VirtualCiYaml: &ci_yaml.VirtualCiYamlFile{},
				Configuration: &cli.Configuration{},
			},
			failSeverity:                      checks.SeverityWarning,
			expectedExit:                      true,
			expectedReportLineCountWithHeader: 3,
			checks: []checks.Check{
				&DummyCheck{name: "DummyCheck3", findings: []checks.CheckFinding{
					{Severity: checks.SeverityWarning},
					{Severity: checks.SeverityWarning},
					{Severity: checks.SeverityError},
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
			err := gcv.SetupFormatter(writer, "table")
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

			lines := strings.Split(strings.TrimSpace(writer.String()), "\n")
			if tt.expectedReportLineCountWithHeader != len(lines) {
				t.Errorf("expected %d lines, got %d", tt.expectedReportLineCountWithHeader, len(lines))
			}
		})
	}
}

func TestCreateCheckInput(t *testing.T) {
	tests := []struct {
		name           string
		configuration  *cli.Configuration
		gitlabCiFile   string
		expectedError  bool
		expectedResult *checks.CheckInput
	}{
		{
			name: "valid configuration with file",
			configuration: &cli.Configuration{
				GitLabCiFile:      "test_data/.gitlab-ci.yml",
				NoLintAPICallInCi: true,
			},
			gitlabCiFile:  "test_data/.gitlab-ci.yml",
			expectedError: false,
			expectedResult: &checks.CheckInput{
				Configuration: &cli.Configuration{
					GitLabCiFile: "test_data/.gitlab-ci.yml",
				},
			},
		},
		{
			name: "invalid configuration with non-existent file",
			configuration: &cli.Configuration{
				GitLabCiFile:      "test_data/non-existent.yml",
				NoLintAPICallInCi: true,
			},
			gitlabCiFile:  "test_data/non-existent.yml",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.CopyFS("test_data/.git", os.DirFS("test_data/.git_template"))
			if err != nil {
				t.Fatal(err)
			}
			ciOldVal := os.Getenv("CI")
			os.Setenv("CI", "true")
			defer func() {
				os.Setenv("CI", ciOldVal)
				err := os.RemoveAll("test_data/.git")
				if err != nil {
					t.Fatal(err)
				}
			}()
			gcv := &GitlabCIVerifier{
				configuration: tt.configuration,
				projectRoot:   "/tmp",
			}

			if tt.gitlabCiFile != "" {
				os.Setenv("GITLAB_CI_FILE", tt.gitlabCiFile)
			} else {
				os.Unsetenv("GITLAB_CI_FILE")
			}

			result, err := gcv.CreateCheckInput()
			if (err != nil) != tt.expectedError {
				t.Errorf("CreateCheckInput() error = %v, wantErr %v", err, tt.expectedError)
				return
			}

			if !tt.expectedError && result.Configuration.GitLabCiFile != tt.expectedResult.Configuration.GitLabCiFile {
				t.Errorf("CreateCheckInput() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestSetupRegoNoPolicies(t *testing.T) {
	gcv := NewGitlabCIVerifier(&cli.Configuration{}, "/tmp")
	gcv.SetupRego()
}
