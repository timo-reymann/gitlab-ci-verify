package shellcheck

import (
	"regexp"
	"testing"
)

func TestValidScript(t *testing.T) {
	sc, err := NewShellChecker()
	if err != nil {
		t.Fatal(err)
	}
	res, err := sc.AnalyzeFile("testdata/valid_script.sh")
	if err != nil {
		t.Fatal(err)
	}

	if res.ExitCode != 0 {
		t.Fatalf("Exptected zero exit code but got %d", res.ExitCode)
	}
}

func TestInValidScript(t *testing.T) {
	sc, err := NewShellChecker()
	if err != nil {
		t.Fatal(err)
	}
	res, err := sc.AnalyzeFile("testdata/invalid_script.sh")
	if err != nil {
		t.Fatal(err)
	}

	if res.ExitCode != 1 {
		t.Fatalf("Exptected failure exit code but got %d", res.ExitCode)
	}
}

func TestShellChecker_Version(t *testing.T) {
	sc, err := NewShellChecker()
	if err != nil {
		t.Fatal(err)
	}

	if match, _ := regexp.Match("[0-9]*\\.[0-9]*\\.[0-9]*", []byte(sc.Version())); !match {
		t.Fatalf("Expted version %s to match semantic release number", sc.Version())
	}
}

func TestAnalyzeSnippet(t *testing.T) {
	s, err := NewShellChecker()
	if err != nil {
		t.Fatal(err.Error())
	}

	var testCases = []struct {
		name                  string
		snippet               []byte
		expectedFindingsCount int
		wantErr               bool
	}{
		{
			name:                  "Valid snippet",
			snippet:               []byte("#!/bin/bash\necho 'Hello World'\n"),
			expectedFindingsCount: 0,
			wantErr:               false,
		},
		{
			name:                  "Snippet with warning",
			snippet:               []byte("#!/bin/bash\nsh echo 'Hello World'\n"),
			expectedFindingsCount: 0,
			wantErr:               false,
		},
		{
			name:                  "Snippet with warning",
			snippet:               []byte("#!/bin/bash\nsh echo $TEST\n"),
			expectedFindingsCount: 1,
			wantErr:               false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := s.AnalyzeSnippet(tc.snippet)

			if (err != nil) != tc.wantErr {
				t.Errorf("Unexpected error: %v", err)
			}

			if tc.expectedFindingsCount != len(result.Findings) {
				t.Errorf("Expected status code: %d, got: %d", tc.expectedFindingsCount, len(result.Findings))
			}
		})
	}
}
