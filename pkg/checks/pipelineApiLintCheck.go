package checks

import (
	"github.com/timo-reymann/gitlab-ci-verify/internal/logging"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/git"
	ci_yaml "github.com/timo-reymann/gitlab-ci-verify/pkg/gitlab/ci-yaml"
	"os"
	"time"
)

type PipelineLintApiCheck struct {
}

func (p PipelineLintApiCheck) Run(i *CheckInput) ([]CheckFinding, error) {
	logging.Verbose("get current working directory")
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	logging.Verbose("get remote urls")
	remoteUrls, err := git.GetRemoteUrls(pwd)
	if err != nil {
		return nil, err
	}

	logging.Verbose("parse remote url contents")
	remoteInfos := git.FilterGitlabRemoteUrls(remoteUrls)

	logging.Verbose("validate ci file against gitlab api")
	res, err := ci_yaml.GetFirstValidationResult(remoteInfos, i.Configuration.GitlabToken, i.Configuration.GitlabBaseUrlOverwrite(), i.CiYaml.FileContent, 3*time.Second)
	if err != nil {
		return nil, err
	}

	if res.LintResult.Valid {
		return []CheckFinding{}, nil
	}

	findings := make([]CheckFinding, 0)

	for _, e := range res.LintResult.Errors {
		findings = append(findings, CheckFinding{
			Severity: SeverityError,
			Code:     "GL-1",
			Line:     -1,
			Message:  e,
			Link:     "https://docs.gitlab.com/ee/ci/yaml",
		})
	}

	for _, w := range res.LintResult.Warnings {
		findings = append(findings, CheckFinding{
			Severity: SeverityWarning,
			Code:     "GL-2",
			Line:     -1,
			Message:  w,
			Link:     "https://docs.gitlab.com/ee/ci/yaml",
		})
	}

	return findings, nil
}
