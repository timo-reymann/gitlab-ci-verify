package formatter

import (
	"fmt"
	"github.com/Ladicle/tabwriter"
	"github.com/fatih/color"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/checks"
	"io"
)

type TextFindingsFormatter struct {
	tabWriter    *tabwriter.Writer
	headingColor *color.Color
}

func (t *TextFindingsFormatter) Init(w io.Writer) error {
	t.tabWriter = tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.TabIndent)
	t.headingColor = color.New(color.Bold)
	return nil
}

func (t *TextFindingsFormatter) Start() error {
	return nil
}

func (t *TextFindingsFormatter) printEntry(heading, val string) error {
	_, err := fmt.Fprintf(t.tabWriter, "%s\t%s\n", t.headingColor.Sprint(heading), val)
	return err
}

func (t *TextFindingsFormatter) Print(finding *checks.CheckFinding) error {
	location, err := finding.Location()
	if err != nil {
		return err
	}

	if err := t.printEntry("Code", finding.Code); err != nil {
		return err
	}

	if err := t.printEntry("Description", finding.Message); err != nil {
		return err
	}

	if err := t.printEntry("Severity", formatSeverity(finding)); err != nil {
		return err
	}

	if err := t.printEntry("Location", "at "+location.String()); err != nil {
		return err
	}

	if err := t.printEntry("Link", finding.Link); err != nil {
		return err
	}

	if err := t.tabWriter.Flush(); err != nil {
		return err
	}

	if _, err = t.tabWriter.Write([]byte("\n")); err != nil {
		return err
	}

	return err
}

func (t *TextFindingsFormatter) End() error {
	return nil
}
