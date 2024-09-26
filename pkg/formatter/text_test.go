package formatter

import (
	"testing"

	"github.com/timo-reymann/gitlab-ci-verify/pkg/checks"
)

func TestTextFindingsFormatter(t *testing.T) {
	for _, tc := range []struct {
		name           string
		findings       []*checks.CheckFinding
		expectedOutput []byte
	}{
		{
			name: "single info finding",
			findings: []*checks.CheckFinding{
				{
					Severity: checks.SeverityInfo,
					Code:     "1",
					Line:     1,
					Message:  "test message goes here",
					Link:     "https://check.link/code",
					File:     "/test.yml",
				},
			},
			expectedOutput: []byte(
				"Code         1\n" +
					"Description  test message goes here\n" +
					"Severity     INFO\n" +
					"Location     /test.yml:1\n" +
					"Link         https://check.link/code\n\n",
			),
		},
		{
			name: "multiple finding",
			findings: []*checks.CheckFinding{
				{
					Severity: checks.SeverityInfo,
					Code:     "1",
					Line:     1,
					Message:  "test message goes here",
					Link:     "https://check.link/code",
					File:     "/test.yml",
				},
				{
					Severity: checks.SeverityStyle,
					Code:     "1",
					Line:     1,
					Message:  "test message goes here",
					Link:     "https://check.link/code",
					File:     "/test.yml",
				},
			},
			expectedOutput: []byte(
				"Code         1\n" +
					"Description  test message goes here\n" +
					"Severity     INFO\n" +
					"Location     /test.yml:1\n" +
					"Link         https://check.link/code\n\n" +
					"Code         1\n" +
					"Description  test message goes here\n" +
					"Severity     STYLE\n" +
					"Location     /test.yml:1\n" +
					"Link         https://check.link/code\n\n",
			),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			verifyFormatter(t, &TextFindingsFormatter{}, tc.findings, tc.expectedOutput)
		})
	}
}
