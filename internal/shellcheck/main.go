package shellcheck

import (
	"github.com/amenzhinsky/go-memexec"
	"github.com/google/shlex"
	"os"
	"strings"
)

// ShellChecker provides API access to the bundled shellcheck binary
type ShellChecker struct {
	exec *memexec.Exec
}

// Version of shellcheck bundled with the application
func (s *ShellChecker) Version() string {
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

	return "N/A"
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
func (s *ShellChecker) AnalyzeFile(path string, extraFlags string) (*Result, error) {
	args := []string{
		"-f", "json",
		"-s", "bash",
	}

	args = append(args, ignoredChecksFlags...)

	if extraFlags != "" {
		extraArgs, err := shlex.Split(extraFlags)
		if err != nil {
			return nil, err
		}
		args = append(args, extraArgs...)
	}

	args = append(args, path)

	return s.execute(args...)
}

// AnalyzeSnippet writes the snippet to a temporary file, analyzes it with shellcheck, and returns the result
func (s *ShellChecker) AnalyzeSnippet(snippet []byte, extraFlags string) (*Result, error) {
	tmpFile, err := os.CreateTemp("", "shellcheck-snippet.*.sh")
	if err != nil {
		return nil, err
	}

	_, err = tmpFile.Write(snippet)
	if err != nil {
		return nil, err
	}

	err = tmpFile.Sync()
	if err != nil {
		return nil, err
	}

	return s.AnalyzeFile(tmpFile.Name(), extraFlags)
}

// NewShellChecker instantiates a new shellcheck instance and loads the binary
func NewShellChecker() (*ShellChecker, error) {
	exec, err := memexec.New(shellcheckBinary)
	if err != nil {
		return nil, err
	}

	return &ShellChecker{exec}, nil
}
