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
		// Determine severity based on warning code
		severity := SeverityWarning
		message := warning.Message

		switch warning.Code {
		case 101:
			// Glob pattern with no matches - Warning level
			message = fmt.Sprintf("Include pattern '%s' did not match any files", warning.IncludePath)
		case 102:
			// Failed to load include file - Error level
			severity = SeverityError
			message = fmt.Sprintf("Include file '%s' could not be loaded: %s", warning.IncludePath, warning.Message)
		}

		finding := l.createFinding(
			i.VirtualCiYaml.EntryFilePath,
			severity,
			warning.Code,
			-1,
			message,
		)
		findings = append(findings, finding)
	}

	return findings, nil
}
