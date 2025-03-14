package checks

import (
	"cmp"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/location"
)

var SeverityError = 0
var SeverityWarning = 1
var SeverityInfo = 2
var SeverityStyle = 3

type CheckFinding struct {
	Severity int
	Code     string
	Line     int
	Message  string
	Link     string
	File     string
}

func (cf *CheckFinding) SeverityName() string {
	return SeverityLevelToName(cf.Severity)
}

func (cf *CheckFinding) Location() (*location.Location, error) {
	loc := location.NewLocation(cf.File, cf.Line)
	return loc.Absolute()
}

func (cf *CheckFinding) Compare(o CheckFinding) int {
	if cf.Severity != o.Severity {
		return cmp.Compare(cf.Severity, o.Severity)
	}

	return cmp.Compare(cf.Line, o.Line)
}
