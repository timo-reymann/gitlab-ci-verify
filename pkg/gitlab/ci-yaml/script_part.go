package ci_yaml

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"strings"
)

// ScriptPart represents a part of a script
type ScriptPart struct {
	// Content is the script content lines
	Content string
	// Node the script was extracted from
	Node *yaml.Node
}

// ScriptPartLine represents a line of a script part
type ScriptPartLine struct {
	// LineContent is a line of a script part
	LineContent string
	// LineNumber contains the number in which the line was found
	LineNumber int
	// Node the script was extracted from
	Node *yaml.Node
}

// Concat all given ScriptPart's into one string and return the lines as well as the buffer
func Concat(parts []ScriptPart) ([]ScriptPartLine, []byte) {
	partLines := make([]ScriptPartLine, 0)
	fullScript := bytes.NewBuffer([]byte{})
	lineNumber := 1
	for _, part := range parts {
		for _, line := range part.SplitContentLines() {
			line.LineNumber = lineNumber
			partLines = append(partLines, line)
			fullScript.Write([]byte(line.LineContent))
			fullScript.Write([]byte("\n"))

			lineNumber++
		}
	}
	return partLines, fullScript.Bytes()
}

// SplitContentLines from a part
func (s *ScriptPart) SplitContentLines() []ScriptPartLine {
	lines := strings.Split(s.Content, "\n")
	lineNo := 1
	partLines := make([]ScriptPartLine, len(lines))
	for idx, line := range lines {
		partLines[idx] = ScriptPartLine{
			LineContent: line,
			LineNumber:  lineNo,
			Node:        s.Node,
		}
		lineNo++
	}
	return partLines
}

func newScriptPart(n *yaml.Node) ScriptPart {
	return ScriptPart{
		Content: n.Value,
		Node:    n,
	}
}
