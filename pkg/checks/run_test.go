package checks

import (
	"errors"
	ciyaml "github.com/timo-reymann/gitlab-ci-verify/pkg/gitlab/ci-yaml"
	"testing"
)

type mockCheck struct {
	Findings []CheckFinding
	Err      error
}

func (m mockCheck) Run(_ *CheckInput) ([]CheckFinding, error) {
	return m.Findings, m.Err
}

func TestRunChecksInParallel(t *testing.T) {
	testCases := []struct {
		name                 string
		checks               []Check
		expectErr            bool
		expectedFindingCount int
	}{
		{
			name:   "no checks",
			checks: []Check{},
		},
		{
			name: "check with errors",
			checks: []Check{
				mockCheck{
					Err: errors.New("failure"),
				},
			},
			expectErr: true,
		},
		{
			name: "check with no errors",
			checks: []Check{
				mockCheck{
					Err: nil,
				},
			},
		},
		{
			name: "check with single result",
			checks: []Check{
				mockCheck{
					Findings: []CheckFinding{
						{
							Severity: SeverityWarning,
							Code:     "123",
							Line:     -1,
							Message:  "test",
						},
					},
				},
			},
			expectedFindingCount: 1,
		},
		{
			name: "check with multiple result",
			checks: []Check{
				mockCheck{
					Findings: []CheckFinding{
						{
							Severity: SeverityWarning,
							Code:     "123",
							Line:     -1,
							Message:  "test",
						},
						{
							Severity: SeverityWarning,
							Code:     "456",
							Line:     -1,
							Message:  "test",
						},
					},
				},
			},
			expectedFindingCount: 2,
		},
		{
			name: "multiple checks with multiple result",
			checks: []Check{
				mockCheck{
					Findings: []CheckFinding{
						{
							Severity: SeverityWarning,
							Code:     "123",
							Line:     -1,
							Message:  "test",
						},
						{
							Severity: SeverityWarning,
							Code:     "456",
							Line:     -1,
							Message:  "test",
						},
					},
				},
				mockCheck{
					Findings: []CheckFinding{
						{
							Severity: SeverityWarning,
							Code:     "789",
							Line:     -1,
							Message:  "test",
						},
						{
							Severity: SeverityWarning,
							Code:     "102",
							Line:     -1,
							Message:  "test",
						},
					},
				},
			},
			expectedFindingCount: 4,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ciYaml, _ := ciyaml.NewCiYamlFile([]byte(``))
			findingsChan := RunChecksInParallel(tc.checks, CheckInput{
				CiYaml: ciYaml,
			}, func(err error) {
				if tc.expectErr && err == nil {
					t.Fatal("Expected error but got none")
				} else if !tc.expectErr && err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}
			})

			findingsCount := 0
			for f := range findingsChan {
				findingsCount += len(f)
			}

			if findingsCount != tc.expectedFindingCount {
				t.Fatalf("Expected %d findings, but got %d", tc.expectedFindingCount, findingsCount)
			}
		})
	}
}
