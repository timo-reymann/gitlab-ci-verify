package checks

import (
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/internal/gitlab/api"
	"github.com/timo-reymann/gitlab-ci-verify/internal/logging"
)

type PipelineLintApiCheck struct {
}

func (p PipelineLintApiCheck) createFinding(path string, severity int, code int, line int, message string) CheckFinding {
	return CheckFinding{
		Severity: severity,
		Code:     fmt.Sprintf("GL-%d", code),
		Line:     line,
		Message:  message,
		Link:     "https://docs.gitlab.com/ee/ci/yaml",
		File:     path,
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
	if !i.HasLintAPIResult() {
		return []CheckFinding{}, nil
	}

	logging.Verbose("validate ci file against gitlab api")
	res := i.LintAPIResult
	if res.LintResult.Valid {
		return []CheckFinding{}, nil
	}

	msgCount := len(res.LintResult.Errors) + len(res.LintResult.Warnings)
	findings := make([]CheckFinding, msgCount)
	idx := 0

	for _, e := range res.LintResult.Errors {
		findings[idx] = p.createFinding(i.VirtualCiYaml.EntryFilePath, SeverityError, 101, -1, p.formatMsg(e))
		idx++
	}

	for _, w := range res.LintResult.Warnings {
		findings[idx] = p.createFinding(i.VirtualCiYaml.EntryFilePath, SeverityWarning, 102, -1, p.formatMsg(w))
		idx++
	}

	if len(findings) == 0 && !res.LintResult.Valid {
		findings = append(findings, p.createFinding(i.VirtualCiYaml.EntryFilePath, SeverityError, 103, -1, "Pipeline is invalid"))
	}

	return findings, nil
}
