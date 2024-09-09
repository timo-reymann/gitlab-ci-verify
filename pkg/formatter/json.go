package formatter

import (
	"encoding/json"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/checks"
	"io"
)

type JsonFindingsFormatter struct {
	writer           io.Writer
	firstItemWritten bool
}

type jsonFinding struct {
	Severity string `json:"severity"`
	Code     string `json:"code"`
	Line     int    `json:"line"`
	Message  string `json:"message"`
	Link     string `json:"link"`
	File     string `json:"file"`
}

func newJsonFinding(f *checks.CheckFinding) jsonFinding {
	return jsonFinding{
		Severity: f.SeverityName(),
		Code:     f.Code,
		Line:     f.Line,
		Message:  f.Message,
		Link:     f.Link,
		File:     f.File,
	}
}

func (j *JsonFindingsFormatter) writeString(val string) error {
	_, err := j.writer.Write([]byte(val))
	return err
}

func (j *JsonFindingsFormatter) Init(w io.Writer) error {
	j.writer = w
	j.firstItemWritten = false
	return nil
}

func (j *JsonFindingsFormatter) Start() error {
	return j.writeString("[")
}

func (j *JsonFindingsFormatter) Print(f *checks.CheckFinding) error {
	buf, err := json.Marshal(newJsonFinding(f))
	if err != nil {
		return err
	}

	if j.firstItemWritten {
		if err := j.writeString(","); err != nil {
			return err
		}
	}

	if err := j.writeString("\n"); err != nil {
		return err
	}

	if err := j.writeString("  "); err != nil {
		return err
	}

	_, err = j.writer.Write(buf)
	j.firstItemWritten = true

	return err
}

func (j *JsonFindingsFormatter) End() error {
	return j.writeString("\n]")
}
