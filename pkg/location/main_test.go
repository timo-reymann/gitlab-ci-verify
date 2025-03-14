package location

import (
	"path/filepath"
	"testing"
)

func TestNewLocation(t *testing.T) {
	testCases := []struct {
		file string
		line int
	}{
		{file: "test.yml", line: 10},
		{file: "/path/to/file.yml", line: 20},
	}

	for _, tc := range testCases {
		t.Run(tc.file, func(t *testing.T) {
			loc := NewLocation(tc.file, tc.line)
			if loc.File != tc.file || loc.Line != tc.line {
				t.Errorf("NewLocation(%q, %d) = %v, want %v", tc.file, tc.line, loc, &Location{File: tc.file, Line: tc.line})
			}
		})
	}
}

func TestLocation_Absolute(t *testing.T) {
	testCases := []struct {
		file    string
		line    int
		wantErr bool
	}{
		{file: "test.yml", line: 10, wantErr: false},
	}

	for _, tc := range testCases {
		t.Run(tc.file, func(t *testing.T) {
			loc := NewLocation(tc.file, tc.line)
			absLoc, err := loc.Absolute()
			if (err != nil) != tc.wantErr {
				t.Errorf("Absolute() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !tc.wantErr {
				expectedPath, _ := filepath.Abs(tc.file)
				if absLoc.File != expectedPath || absLoc.Line != tc.line {
					t.Errorf("Absolute() = %v, want %v", absLoc, &Location{File: expectedPath, Line: tc.line})
				}
			}
		})
	}
}

func TestLocation_String(t *testing.T) {
	testCases := []struct {
		file string
		line int
		want string
	}{
		{file: "test.yml", line: 10, want: "test.yml:10"},
		{file: "/path/to/file.yml", line: 20, want: "/path/to/file.yml:20"},
	}

	for _, tc := range testCases {
		t.Run(tc.file, func(t *testing.T) {
			loc := NewLocation(tc.file, tc.line)
			got := loc.String()
			if got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}
