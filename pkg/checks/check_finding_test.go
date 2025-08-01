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
	if loc.File != "/test.yml" && loc.Line != 20 {
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

func TestHasEqualOrHigherSeverityThan(t *testing.T) {
	tests := []struct {
		name                 string
		findingSeverity      int
		checkAgainstSeverity int
		expected             bool
	}{
		{
			name:                 "equal severities",
			findingSeverity:      SeverityError,
			checkAgainstSeverity: SeverityError,
			expected:             true,
		},
		{
			name:                 "higher severity",
			findingSeverity:      SeverityError,
			checkAgainstSeverity: SeverityWarning,
			expected:             true,
		},
		{
			name:                 "lower severity",
			findingSeverity:      SeverityWarning,
			checkAgainstSeverity: SeverityError,
			expected:             false,
		},
		{
			name:                 "equal lowest severities",
			findingSeverity:      SeverityInfo,
			checkAgainstSeverity: SeverityInfo,
			expected:             true,
		},
		{
			name:                 "higher than lowest severity",
			findingSeverity:      SeverityWarning,
			checkAgainstSeverity: SeverityInfo,
			expected:             true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := CheckFinding{Severity: tt.findingSeverity}
			result := f.HasEqualOrHigherSeverityThan(tt.checkAgainstSeverity)
			if result != tt.expected {
				t.Errorf("HasEqualOrHigherSeverityThan(%v, %v) = %v, want %v", tt.findingSeverity, tt.checkAgainstSeverity, result, tt.expected)
			}
		})
	}
}

func TestCheckFinding_Fingerprint(t *testing.T) {
	tests := []struct {
		name    string
		finding CheckFinding
		want    string
	}{
		{
			name: "basic finding",
			finding: CheckFinding{
				Severity: 1,
				Code:     "TEST-001",
				Line:     10,
				Message:  "Test message",
				Link:     "https://example.com",
				File:     "test.yml",
			},
			want: "ea9d9494671bbc9455169bd1c4c71eb2aab34689a41a6bbb9f908401aa8d3e9b",
		},
		{
			name: "empty fields",
			finding: CheckFinding{
				Severity: 0,
				Code:     "",
				Line:     0,
				Message:  "",
				Link:     "",
				File:     "",
			},
			want: "f1534392279bddbf9d43dde8701cb5be14b82f76ec6607bf8d6ad557f60f304e",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.finding.Fingerprint(); got != tt.want {
				t.Errorf("CheckFinding.Fingerprint() = %v, want %v", got, tt.want)
			}

			// Test consistency - multiple calls should return the same fingerprint
			if second := tt.finding.Fingerprint(); second != tt.want {
				t.Errorf("CheckFinding.Fingerprint() second call = %v, want %v", second, tt.want)
			}
		})
	}
}
