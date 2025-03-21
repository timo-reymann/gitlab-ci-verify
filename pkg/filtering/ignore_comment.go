package filtering

import (
	"gopkg.in/yaml.v3"
	"strings"
)

// IgnoreComment represents a comment for codes that should be ignored
// in the verification process. An YAML comment might result in multiple
// IgnoreComment structs if multiple codes are present in the same comment
type IgnoreComment struct {
	// Comment is the full comment that was found
	Comment string
	// Code is the code that should be ignored
	Code string
}

// IgnoreCommentsToCodes converts a list of IgnoreComment structs to a list of codes
func IgnoreCommentsToCodes(comments []IgnoreComment) []string {
	codes := make([]string, 0)
	for _, comment := range comments {
		codes = append(codes, comment.Code)
	}
	return codes
}

// ParseIgnoreForLine parses the ignored codes for a given line
// and returns a list of IgnoreComment structs for each found ignore code
func ParseIgnoreForLine(lineNumberMapping yaml.LineNumberMapping, lines []int) []IgnoreComment {
	ignoreComments := make([]IgnoreComment, 0)

	for _, line := range lines {
		nodes, ok := lineNumberMapping[line]
		if !ok {
			continue
		}

		for _, node := range nodes {
			if lineComment := ParseIgnoreComment(node.LineComment); lineComment != nil {
				ignoreComments = append(ignoreComments, lineComment...)
			}

			if blockComment := ParseIgnoreComment(node.HeadComment); blockComment != nil {
				ignoreComments = append(ignoreComments, blockComment...)
			}
		}
	}

	return ignoreComments
}

// ParseIgnoreComment parses the ignored codes from a given comment
// creating a IgnoreComment struct for each found ignored code in the comment
func ParseIgnoreComment(comment string) []IgnoreComment {
	parsed := make([]IgnoreComment, 0)

	if !strings.Contains(comment, "gitlab-ci-verify") {
		return nil
	}

	// Process multi line comments
	if strings.Contains(comment, "\n") {
		lines := strings.Split(comment, "\n")
		for _, line := range lines {
			parsed = append(parsed, ParseIgnoreComment(line)...)
		}
		return parsed
	}

	parts := strings.Split(comment, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, "ignore:") {
			code := strings.TrimPrefix(part, "ignore:")
			parsed = append(parsed, IgnoreComment{
				Comment: strings.TrimSpace(strings.TrimPrefix(comment, "#")),
				Code:    code,
			})
		}
	}

	return parsed
}
