package checks

import (
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/internal/logging"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/git"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/gitlab/api"
	ci_yaml "github.com/timo-reymann/gitlab-ci-verify/pkg/gitlab/ci-yaml"
	"os"
	"time"
)

type PipelineLintApiCheck struct {
}

func (p PipelineLintApiCheck) createFinding(severity int, code int, line int, message string) CheckFinding {
	return CheckFinding{
		Severity: severity,
		Code:     fmt.Sprintf("GL-%d", code),
		Line:     line,
		Message:  message,
		Link:     "https://docs.gitlab.com/ee/ci/yaml",
	}
}

func (p PipelineLintApiCheck) formatMsg(msg string) string {
	job, msg := api.ParsePipelineMessage(msg)
	if job == "" {
		return msg
	}
	return fmt.Sprintf("[%s] %s", job, msg)
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

	msgCount := len(res.LintResult.Errors) + len(res.LintResult.Warnings)
	findings := make([]CheckFinding, msgCount)
	idx := 0

	for _, e := range res.LintResult.Errors {
		findings[idx] = p.createFinding(SeverityError, 101, -1, p.formatMsg(e))
		idx++
	}

	for _, w := range res.LintResult.Warnings {
		findings[idx] = p.createFinding(SeverityWarning, 102, -1, p.formatMsg(w))
		idx++
	}

	if len(findings) == 0 {
		findings = append(findings, p.createFinding(SeverityError, 103, -1, "Pipeline is invalid"))
	}

	return findings, nil
}
