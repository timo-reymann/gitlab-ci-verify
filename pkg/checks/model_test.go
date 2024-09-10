package checks

import (
	"fmt"
	"testing"
)

func TestCheckFinding_SeverityName(t *testing.T) {
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

func TestCheckFinding_Location(t *testing.T) {
	cf := CheckFinding{File: "/test.yml", Line: 20}
	loc, err := cf.Location()
	if err != nil {
		t.Fatal(err)
	}
	if loc != "/test.yml:20" {
		t.Fatal("Location generation does not work")
	}
}

func TestCheckFinding_Compare(t *testing.T) {
	testCases := []struct {
		name     string
		cf       *CheckFinding
		other    CheckFinding
		expected int
	}{
		{
			name:     "Valid 1 (Equal Severity, Equal Line)",
			cf:       &CheckFinding{Severity: SeverityError, Line: 10},
			other:    CheckFinding{Severity: SeverityError, Line: 10},
			expected: 0,
		},
		{
			name:     "Valid 2 (Equal Severity, Different Line)",
			cf:       &CheckFinding{Severity: SeverityError, Line: 10},
			other:    CheckFinding{Severity: SeverityError, Line: 20},
			expected: -1,
		},
		{
			name:     "Valid 3 (Different Severity, Equal Line)",
			cf:       &CheckFinding{Severity: SeverityWarning, Line: 10},
			other:    CheckFinding{Severity: SeverityError, Line: 10},
			expected: 1,
		},
		{
			name:     "Valid 4 (Different Severity, Different Line)",
			cf:       &CheckFinding{Severity: SeverityError, Line: 10},
			other:    CheckFinding{Severity: SeverityInfo, Line: 20},
			expected: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.cf.Compare(tc.other)
			if actual != tc.expected {
				t.Errorf("Test case %s: Expected %d, Got %d", tc.name, tc.expected, actual)
			}
		})
	}
}
