package checks

import (
	"github.com/timo-reymann/gitlab-ci-verify/pkg/cli"
	"path"
	"testing"
)

func TestRegoBundleCheck_Run(t *testing.T) {
	c := RegoBundleCheck{
		BundlePath: "test_data/rego-bundle",
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
			verifyFindings(t, tc.expectedFindings, checkMustSucceed(c.Run(&CheckInput{
				CiYaml:        newCiYamlFromFile(t, path.Join("test_data", "ci_yamls", tc.file)),
				MergedCiYaml:  newCiYamlFromFile(t, path.Join("test_data", "ci_yamls", tc.file)),
				Configuration: &cli.Configuration{},
			})))
		})
	}
}
