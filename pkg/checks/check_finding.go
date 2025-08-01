package checks

import (
	"cmp"
	"crypto/sha256"
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/v2/pkg/location"
	"slices"
)

// CheckFinding represents a finding from a check
type CheckFinding struct {
	// Severity of the finding
	Severity int
	// Code of the finding, unique identifier in format <CHECK SHORT>-<CODE>
	Code string
	// Line of the finding
	Line int
	// Message of the finding
	Message string
	// Link to the documentation for more information
	Link string
	// File path where the finding was found
	File string
}

// SeverityLevelToName converts a severity level to a human readable name
func (cf *CheckFinding) SeverityName() string {
	return SeverityLevelToName(cf.Severity)
}

// HasEqualOrHigherSeverityThan checks if the finding has severity equal or higher than the given severity
func (cf *CheckFinding) HasEqualOrHigherSeverityThan(severity int) bool {
	return cf.Severity <= severity
}

// HasCodeIn checks if the finding has a code in the given list of codes
func (cf *CheckFinding) HasCodeIn(codes []string) bool {
	return slices.Contains(codes, cf.Code)
}

// Location returns the location of the finding as an absolute location
func (cf *CheckFinding) Location() (*location.Location, error) {
	loc := location.NewLocation(cf.File, cf.Line)
	return loc.Absolute()
}

// Compare two findings are equal
func (cf *CheckFinding) Compare(o CheckFinding) int {
	if cf.Severity != o.Severity {
		return cmp.Compare(cf.Severity, o.Severity)
	}

	return cmp.Compare(cf.Line, o.Line)
}

func (cf *CheckFinding) Fingerprint() string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%d%s%d%s%s%s", cf.Severity, cf.Code, cf.Line, cf.Message, cf.Link, cf.File)))
	return fmt.Sprintf("%x", h.Sum(nil))
}
