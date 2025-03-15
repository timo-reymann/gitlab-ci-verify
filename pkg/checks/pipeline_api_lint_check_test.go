package checks

import (
	"encoding/json"
	"github.com/timo-reymann/gitlab-ci-verify/internal/cli"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/git"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/gitlab/api"
	ciyaml "github.com/timo-reymann/gitlab-ci-verify/pkg/gitlab/ci-yaml"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func mockCiValidate(lintResult api.CiLintResult) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		jsonRes, err := json.Marshal(lintResult)
		if err != nil {
			panic(err)
		}
		_, _ = writer.Write(jsonRes)
	}))
}

func TestPipelineLintApiCheck_Run(t *testing.T) {
	testCases := []struct {
		name             string
		folder           string
		lintResult       api.CiLintResult
		expectedFindings []CheckFinding
		ciEnvVarVal      string
	}{
		{
			name:             "Test project with empty ci yaml and valid pipeline",
			folder:           "project_with_git_repo_empty_ci_yaml.git",
			expectedFindings: []CheckFinding{},
			lintResult: api.CiLintResult{
				Valid: true,
			},
		},
		{
			name:   "Test project with empty ci yaml and invalid pipeline",
			folder: "project_with_git_repo_empty_ci_yaml.git",
			lintResult: api.CiLintResult{
				Valid: false,
			},
			expectedFindings: []CheckFinding{
				{
					Severity: SeverityError,
					Code:     "GL-103",
					Line:     -1,
					Message:  "Pipeline is invalid",
					Link:     "https://docs.gitlab.com/ee/ci/yaml",
					File:     ".gitlab-ci.yml",
				},
			},
		},
		{
			name:   "Test project with empty ci yaml and invalid pipeline with warnings",
			folder: "project_with_git_repo_empty_ci_yaml.git",
			lintResult: api.CiLintResult{
				Valid:    false,
				Warnings: []string{"this bad"},
			},
			expectedFindings: []CheckFinding{
				{
					Severity: SeverityWarning,
					Code:     "GL-102",
					Line:     -1,
					Message:  "this bad",
					Link:     "https://docs.gitlab.com/ee/ci/yaml",
					File:     ".gitlab-ci.yml",
				},
			},
		},
		{
			name:   "Test project with empty ci yaml and invalid pipeline with errors",
			folder: "project_with_git_repo_empty_ci_yaml.git",
			lintResult: api.CiLintResult{
				Valid:  false,
				Errors: []string{"this bad"},
			},
			expectedFindings: []CheckFinding{
				{
					Severity: SeverityError,
					Code:     "GL-101",
					Line:     -1,
					Message:  "this bad",
					Link:     "https://docs.gitlab.com/ee/ci/yaml",
					File:     ".gitlab-ci.yml",
				},
			},
		},
	}
	c := PipelineLintApiCheck{}
	oldCwd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(oldCwd)
	}()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = os.Setenv("CI", tc.ciEnvVarVal)
			projectRoot := path.Join("test_data", tc.folder)
			ciValidateMockServer := mockCiValidate(tc.lintResult)

			_ = os.Chdir(projectRoot)
			ciYaml, err := ciyaml.NewCiYamlFile([]byte(`{}`))
			if err != nil {
				t.Fatal(err)
			}

			checkInput := createCheckInput(t, ciYaml, projectRoot, ".gitlab-ci.yml")
			checkInput.LintAPIResult = &ciyaml.VerificationResultWithRemoteInfo{
				RemoteInfo: &git.GitlabRemoteUrlInfo{
					Hostname:       ciValidateMockServer.URL,
					ClonedViaHttps: true,
				},
				LintResult: &tc.lintResult,
			}
			checkInput.Configuration = &cli.Configuration{
				NoLintAPICallInCi: false,
			}

			VerifyFindings(t, tc.expectedFindings, CheckMustSucceed(c.Run(checkInput)))
		})
		_ = os.Setenv("CI", "")
	}
}
