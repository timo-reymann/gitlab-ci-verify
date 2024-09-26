package formatter

import (
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/checks"
	"io"
)

type FindingsFormatter interface {
	Init(w io.Writer) error
	Start() error
	Print(finding *checks.CheckFinding) error
	End() error
}

func Get(format string) (FindingsFormatter, error) {
	switch format {
	case "table":
		return &TableFindingsFormatter{}, nil
	case "text":
		return &TextFindingsFormatter{}, nil
	case "json":
		return &JsonFindingsFormatter{}, nil
	}
	return nil, fmt.Errorf("unsupported format %s", format)
}
