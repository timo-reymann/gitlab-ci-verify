package buildinfo

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintVersionInfo(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	PrintVersionInfo("gitlab-ci-verify", buf)
	lines := strings.Split(buf.String(), "\n")
	linesLen := len(lines)
	expectedLinesLen := 10
	if linesLen != expectedLinesLen {
		t.Fatalf("Expected %d lines to be printed, but only %d were found", expectedLinesLen, linesLen)
	}
}
