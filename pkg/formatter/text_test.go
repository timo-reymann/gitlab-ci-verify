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
				},
			},
			expectedOutput: []byte(
				"Severity  Code  Line  Description             Link\n" +
					"INFO      1     1     test message goes here  https://check.link/code\n",
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
				},
				{
					Severity: checks.SeverityStyle,
					Code:     "1",
					Line:     1,
					Message:  "test message goes here",
					Link:     "https://check.link/code",
				},
			},
			expectedOutput: []byte(
				"Severity  Code  Line  Description             Link\n" +
					"INFO      1     1     test message goes here  https://check.link/code\n" +
					"STYLE     1     1     test message goes here  https://check.link/code\n",
			),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			verifyFormatter(t, &TextFindingsFormatter{}, tc.findings, tc.expectedOutput)
		})
	}
}
