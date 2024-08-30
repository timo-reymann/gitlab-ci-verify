package shellcheck

import (
	"github.com/amenzhinsky/go-memexec"
	"strings"
)

// ShellChecker provides API access to the bundled shellcheck binary
type ShellChecker struct {
	exec *memexec.Exec
}

// Version of shellcheck bundled with the application
func (s *ShellChecker) Version() string {
	version := "N/A"
	cmd := s.exec.Command("--version")
	output, _ := cmd.Output()
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		if key == "version" {
			return strings.TrimSpace(parts[1])
		}
	}

	return version
}

// Close handle to shellcheck
func (s *ShellChecker) Close() error {
	return s.exec.Close()
}

func (s *ShellChecker) execute(args ...string) (*Result, error) {
	cmd := s.exec.Command(args...)
	output, _ := cmd.Output()
	return NewResult(cmd.ProcessState.ExitCode(), output)
}

// AnalyzeFile for a given path
func (s *ShellChecker) AnalyzeFile(path string) (*Result, error) {
	return s.execute("-f", "json", "-s", "bash", path)
}

// NewShellChecker instantiates a new shellcheck instance and loads the binary
func NewShellChecker() (*ShellChecker, error) {
	exec, err := memexec.New(shellcheckBinary)
	if err != nil {
		return nil, err
	}

	return &ShellChecker{exec}, nil
}
