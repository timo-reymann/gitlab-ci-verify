package formatter

import (
	"github.com/timo-reymann/gitlab-ci-verify/pkg/checks"
	"io"
)

type FindingsFormatter interface {
	Init(w io.Writer) error
	Start() error
	Print(finding *checks.CheckFinding) error
	End() error
}
