package formatter

import (
	"testing"

	"github.com/timo-reymann/gitlab-ci-verify/v2/pkg/checks"
)

func TestJsonFindingsFormatter(t *testing.T) {
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
					File:     ".gitlab-ci.yml",
				},
			},
			expectedOutput: []byte(
				"[\n  {\"severity\":\"Info\",\"code\":\"1\",\"line\":1,\"message\":\"test message goes here\",\"link\":\"https://check.link/code\",\"file\":\".gitlab-ci.yml\"}\n]",
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
					File:     ".gitlab-ci.yml",
				},
				{
					Severity: checks.SeverityStyle,
					Code:     "1",
					Line:     1,
					Message:  "test message goes here",
					Link:     "https://check.link/code",
					File:     ".gitlab-ci-1.yml",
				},
			},
			expectedOutput: []byte(
				"[\n  {\"severity\":\"Info\",\"code\":\"1\",\"line\":1,\"message\":\"test message goes here\",\"link\":\"https://check.link/code\",\"file\":\".gitlab-ci.yml\"},\n  {\"severity\":\"Style\",\"code\":\"1\",\"line\":1,\"message\":\"test message goes here\",\"link\":\"https://check.link/code\",\"file\":\".gitlab-ci-1.yml\"}\n]",
			),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			verifyFormatter(t, &JsonFindingsFormatter{}, tc.findings, tc.expectedOutput)
		})
	}
}
