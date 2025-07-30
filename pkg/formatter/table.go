package formatter

import (
	"fmt"
	"github.com/Ladicle/tabwriter"
	"github.com/fatih/color"
	"github.com/timo-reymann/gitlab-ci-verify/v2/pkg/checks"
	"io"
	"strconv"
	"strings"
)

type TableFindingsFormatter struct {
	tabWriter *tabwriter.Writer
}

func (t *TableFindingsFormatter) Init(w io.Writer) error {
	t.tabWriter = tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.TabIndent)
	return nil
}

func (t *TableFindingsFormatter) Start() error {
	c := color.New(color.Bold)
	_, err := fmt.Fprintln(t.tabWriter, strings.Join([]string{
		c.Sprint("Severity"),
		c.Sprint("Code"),
		c.Sprint("Line"),
		c.Sprint("Description"),
		c.Sprint("Link"),
		c.Sprintf("Location"),
	}, "\t"))
	return err
}

func (t *TableFindingsFormatter) Print(finding *checks.CheckFinding) error {
	location, err := finding.Location()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(t.tabWriter, strings.Join([]string{
		formatSeverity(finding),
		finding.Code,
		strconv.Itoa(finding.Line),
		finding.Message,
		finding.Link,
		location.String(),
	}, "\t"))
	return err
}

func (t *TableFindingsFormatter) End() error {
	return t.tabWriter.Flush()
}
