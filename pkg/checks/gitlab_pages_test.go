package checks

import (
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
					Message:  "Pages job does not contain artifacts with a public path.",
					Link:     "https://gitlab-ci-verify.timo-reymann.de/findings/GL-202.html",
					File:     "emptyArtifactPaths.yml",
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
					Message:  "Pages job does not contain artifacts with a public path.",
					Link:     "https://gitlab-ci-verify.timo-reymann.de/findings/GL-202.html",
					File:     "invalidArtifacts.yml",
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
					Message:  "Pages job does not define artifacts.",
					Link:     "https://gitlab-ci-verify.timo-reymann.de/findings/GL-202.html",
					File:     "noArtifactPaths.yml",
				},
			},
		},
		{
			name: "with included invalid artifacts",
			file: "includingInvalidArtifacts.yml",
			expectedFindings: []CheckFinding{
				{
					Severity: SeverityWarning,
					Code:     "GL-201",
					Line:     3,
					Message:  "Pages job does not contain artifacts with a public path.",
					Link:     "https://gitlab-ci-verify.timo-reymann.de/findings/GL-202.html",
					File:     "test_data/pages/invalidArtifacts.yml",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRoot := path.Join("test_data", "pages")
			input := createCheckInput(t, NewCiYamlFromFile(t, path.Join(projectRoot, tc.file)), projectRoot, tc.file)
			VerifyFindings(t, tc.expectedFindings, CheckMustSucceed(c.Run(input)))
		})
	}
}
