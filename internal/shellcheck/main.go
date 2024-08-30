package shellcheck

import (
	"github.com/amenzhinsky/go-memexec"
)

type ShellChecker struct {
	exec *memexec.Exec
}

func (s *ShellChecker) Close() error {
	return s.exec.Close()
}

func (s *ShellChecker) execute(args ...string) (*Result, error) {
	cmd := s.exec.Command(args...)
	output, _ := cmd.Output()
	return NewResult(cmd.ProcessState.ExitCode(), output)
}

func (s *ShellChecker) AnalyzeFile(path string) (*Result, error) {
	return s.execute("-f", "json", "-s", "bash", path)
}

func NewShellChecker() (*ShellChecker, error) {
	exec, err := memexec.New(shellcheckBinary)
	if err != nil {
		return nil, err
	}

	return &ShellChecker{exec}, nil
}
