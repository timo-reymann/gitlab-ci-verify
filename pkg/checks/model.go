package checks

import (
	"cmp"
	"fmt"
	"path/filepath"
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

func (cf *CheckFinding) Location() (string, error) {
	absPath, err := filepath.Abs(cf.File)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%d", absPath, cf.Line), nil
}

func (cf *CheckFinding) Compare(o CheckFinding) int {
	if cf.Severity != o.Severity {
		return cmp.Compare(cf.Severity, o.Severity)
	}

	return cmp.Compare(cf.Line, o.Line)
}
