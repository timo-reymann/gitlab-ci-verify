package formatter

import (
	"fmt"
	"github.com/Ladicle/tabwriter"
	"github.com/fatih/color"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/checks"
	"io"
	"strconv"
	"strings"
)

type TextFindingsFormatter struct {
	tabWriter *tabwriter.Writer
}

func (t *TextFindingsFormatter) Init(w io.Writer) error {
	t.tabWriter = tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.TabIndent)
	return nil
}

func (t *TextFindingsFormatter) Start() error {
	c := color.New(color.Bold)
	_, err := fmt.Fprintln(t.tabWriter, strings.Join([]string{
		c.Sprint("Severity"),
		c.Sprint("Code"),
		c.Sprint("Line"),
		c.Sprint("Description"),
		c.Sprint("Link"),
	}, "\t"))
	return err
}

func (t *TextFindingsFormatter) severity(cf *checks.CheckFinding) string {
	var colorize func(string, ...any) string
	switch cf.Severity {
	case checks.SeverityError:
		colorize = color.RedString
	case checks.SeverityWarning:
		colorize = color.YellowString
	case checks.SeverityInfo:
		colorize = color.CyanString
	case checks.SeverityStyle:
		colorize = color.BlueString
	}

	return colorize(strings.ToUpper(cf.SeverityName()))
}

func (t *TextFindingsFormatter) Print(finding *checks.CheckFinding) error {
	_, err := fmt.Fprintln(t.tabWriter, strings.Join([]string{
		t.severity(finding),
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
