package shellcheck

import "encoding/json"

type Result struct {
	ExitCode  int
	RawResult []byte
	Findings  []Finding
}

type Finding struct {
	File      string `json:"file"`
	Line      int    `json:"line"`
	EndLine   int    `json:"endLine"`
	Column    int    `json:"column"`
	EndColumn int    `json:"endColumn"`
	Level     string `json:"level"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Fix       Fix    `json:"fix"`
}

type Fix struct {
	Replacements []struct {
		Column         int    `json:"column"`
		EndColumn      int    `json:"endColumn"`
		EndLine        int    `json:"endLine"`
		InsertionPoint string `json:"insertionPoint"`
		Line           int    `json:"line"`
		Precedence     int    `json:"precedence"`
		Replacement    string `json:"replacement"`
	} `json:"replacements"`
}

func (s *Result) parseJSON() ([]Finding, error) {
	parsed := make([]Finding, 0)
	err := json.Unmarshal(s.RawResult, &parsed)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

func NewResult(exitCode int, output []byte) (*Result, error) {
	result := Result{
		ExitCode:  exitCode,
		RawResult: output,
	}
	findings, err := result.parseJSON()
	if err != nil {
		return nil, err
	}

	result.Findings = findings

	return &result, nil
}
