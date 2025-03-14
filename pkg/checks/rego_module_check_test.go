package checks

import (
	"path"
	"testing"
)

func TestRegoModuleCheck_Run(t *testing.T) {
	c := ModuleCheck{
		ModulePath: "test_data/rego-bundle/always_trigger.rego",
	}
	testCases := []struct {
		name             string
		file             string
		expectedFindings []CheckFinding
	}{
		{
			name: "with valid artifacts and job",
			file: "singleBuild.yml",
			expectedFindings: []CheckFinding{
				{
					Severity: SeverityWarning,
					Code:     "420",
					Line:     -1,
					Message:  "always triggers",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := createCheckInput(t, NewCiYamlFromFile(t, path.Join("test_data", "ci_yamls", tc.file)), ".", tc.file)
			VerifyFindings(t, tc.expectedFindings, CheckMustSucceed(c.Run(input)))
		})
	}
}
