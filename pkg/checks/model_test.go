package checks

import (
	"fmt"
	"testing"
)

func TestSeverityName(t *testing.T) {
	testCases := []struct {
		severity int
		want     string
	}{
		{SeverityError, "Error"},
		{SeverityWarning, "Warning"},
		{SeverityInfo, "Info"},
		{SeverityStyle, "Style"},
		{-1, ""}, // Invalid severity
		{4, ""},  // Invalid severity
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("TestSeverityName_%d", tc.severity), func(t *testing.T) {
			cf := CheckFinding{Severity: tc.severity}
			got := cf.SeverityName()
			if got != tc.want {
				t.Errorf("SeverityName(%d) = %q, want %q", tc.severity, got, tc.want)
			}
		})
	}
}
