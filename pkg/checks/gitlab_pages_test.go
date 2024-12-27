package checks

import (
	"github.com/timo-reymann/gitlab-ci-verify/pkg/cli"
	"path"
	"testing"
)

func TestGitlabPagesJobCheck_Run(t *testing.T) {
	c := NewGitlabPagesJobCheck()
	testCases := []struct {
		name             string
		file             string
		expectedFindings []CheckFinding
	}{
		{
			name:             "with valid artifacts and job",
			file:             "validArtifacts.yml",
			expectedFindings: []CheckFinding{},
		},
		{
			name: "with empty artifact paths and job",
			file: "emptyArtifactPaths.yml",
			expectedFindings: []CheckFinding{
				{
					Severity: SeverityWarning,
					Code:     "GL-201",
					Line:     2,
					Message:  "pages job does not contain artifacts with public path",
					Link:     "https://docs.gitlab.com/ee/user/project/pages",
				},
			},
		},
		{
			name: "with invalid artifacts path and job",
			file: "invalidArtifacts.yml",
			expectedFindings: []CheckFinding{
				{
					Severity: SeverityWarning,
					Code:     "GL-201",
					Line:     3,
					Message:  "pages job does not contain artifacts with public path",
					Link:     "https://docs.gitlab.com/ee/user/project/pages",
				},
			},
		},
		{
			name: "with empty artifacts configuration",
			file: "noArtifactPaths.yml",
			expectedFindings: []CheckFinding{
				{
					Severity: SeverityWarning,
					Code:     "GL-202",
					Line:     1,
					Message:  "pages job does not define artifacts",
					Link:     "https://docs.gitlab.com/ee/user/project/pages",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			VerifyFindings(t, tc.expectedFindings, CheckMustSucceed(c.Run(&CheckInput{
				CiYaml:        NewCiYamlFromFile(t, path.Join("test_data", "pages", tc.file)),
				MergedCiYaml:  NewCiYamlFromFile(t, path.Join("test_data", "pages", tc.file)),
				Configuration: &cli.Configuration{},
			})))
		})
	}
}
