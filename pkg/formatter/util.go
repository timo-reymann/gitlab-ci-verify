package formatter

import (
	"github.com/fatih/color"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/checks"
	"strings"
)

func formatSeverity(cf *checks.CheckFinding) string {
	var colorize func(string, ...any) string
	switch cf.Severity {
	case checks.SeverityError:
		colorize = color.RedString
	case checks.SeverityWarning:
		colorize = color.YellowString
	case checks.SeverityInfo:
		colorize = color.CyanString
	case checks.SeverityStyle:
		colorize = color.BlueString
	}

	return colorize(strings.ToUpper(cf.SeverityName()))
}
