package checks

import "fmt"

type LocalIncludeGlobCheck struct {
}

func (l LocalIncludeGlobCheck) createFinding(path string, severity int, code int, line int, message string) CheckFinding {
	findingId := fmt.Sprintf("INC-%d", code)
	return CheckFinding{
		Severity: severity,
		Code:     findingId,
		Line:     line,
		Message:  message,
		Link:     fmt.Sprintf("https://gitlab-ci-verify.timo-reymann.de/findings/%s.html", findingId),
		File:     path,
	}
}

func (l LocalIncludeGlobCheck) Run(i *CheckInput) ([]CheckFinding, error) {
	findings := make([]CheckFinding, 0)

	for _, warning := range i.VirtualCiYaml.Warnings {
		finding := l.createFinding(
			i.VirtualCiYaml.EntryFilePath,
			SeverityWarning,
			101,
			-1,
			fmt.Sprintf("Include pattern '%s' did not match any files", warning.IncludePath),
		)
		findings = append(findings, finding)
	}

	return findings, nil
}
