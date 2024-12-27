package checks

import (
	"github.com/timo-reymann/gitlab-ci-verify/pkg/cli"
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
			VerifyFindings(t, tc.expectedFindings, CheckMustSucceed(c.Run(&CheckInput{
				CiYaml:        NewCiYamlFromFile(t, path.Join("test_data", "ci_yamls", tc.file)),
				MergedCiYaml:  NewCiYamlFromFile(t, path.Join("test_data", "ci_yamls", tc.file)),
				Configuration: &cli.Configuration{},
			})))
		})
	}
}
