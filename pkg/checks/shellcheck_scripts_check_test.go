package checks

import (
	"github.com/timo-reymann/gitlab-ci-verify/pkg/cli"
	"path"
	"testing"
)

func TestShellScriptCheck_Run(t *testing.T) {
	c := ShellScriptCheck{}
	testCases := []struct {
		name             string
		file             string
		expectedFindings []CheckFinding
	}{
		{
			name: "With info finding",
			file: "withScriptInfoFinding.yml",
			expectedFindings: []CheckFinding{
				{
					Severity: SeverityInfo,
					Code:     "SC-2086",
					Line:     4,
					Message:  "[build:script:1] Double quote to prevent globbing and word splitting.",
					Link:     "https://www.shellcheck.net/wiki/SC2086",
				},
			},
		},
		{
			name: "With error finding",
			file: "withScriptErrorFinding.yml",
			expectedFindings: []CheckFinding{
				{
					Severity: SeverityError,
					Code:     "SC-2189",
					Line:     4,
					Message:  "[build:script:1] You can't have | between this redirection and the command it should apply to.",
					Link:     "https://www.shellcheck.net/wiki/SC2189",
				},
			},
		},
		{
			name: "With warning finding",
			file: "withScriptWarningFinding.yml",
			expectedFindings: []CheckFinding{
				{
					Severity: SeverityWarning,
					Code:     "SC-2202",
					Line:     4,
					Message:  "[build:script:1] Globs don't work as operands in [ ]. Use a loop.",
					Link:     "https://www.shellcheck.net/wiki/SC2202",
				},
			},
		},
		{
			name: "With style finding",
			file: "withStyleWarningFinding.yml",
			expectedFindings: []CheckFinding{
				{
					Severity: SeverityStyle,
					Code:     "SC-2006",
					Line:     4,
					Message:  "[build:script:1] Use $(...) notation instead of legacy backticks `...`.",
					Link:     "https://www.shellcheck.net/wiki/SC2006",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			verifyFindings(t, tc.expectedFindings, checkMustSucceed(c.Run(&CheckInput{
				CiYaml:        newCiYamlFromFile(t, path.Join("test_data", "ci_yamls", tc.file)),
				Configuration: &cli.Configuration{},
			})))
		})
	}
}
