package formatter

import (
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/checks"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"
)

type TextFindingsFormatter struct {
	tabWriter *tabwriter.Writer
}

func (t *TextFindingsFormatter) Init(w io.Writer) error {
	t.tabWriter = tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.TabIndent)
	return nil
}

func (t *TextFindingsFormatter) Start() error {
	_, err := fmt.Fprintln(t.tabWriter, strings.Join([]string{
		"Severity",
		"Code",
		"Line",
		"Description",
		"Link",
	}, "\t"))
	return err
}

func (t *TextFindingsFormatter) Print(finding *checks.CheckFinding) error {
	_, err := fmt.Fprintln(t.tabWriter, strings.Join([]string{
		strings.ToUpper(finding.SeverityName()),
		finding.Code,
		strconv.Itoa(finding.Line),
		finding.Message,
		finding.Link,
	}, "\t"))
	return err
}

func (t *TextFindingsFormatter) End() error {
	return t.tabWriter.Flush()
}
