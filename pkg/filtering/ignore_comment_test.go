package filtering

import (
	"gopkg.in/yaml.v3"
	"reflect"
	"testing"
)

func TestParseIgnoreComment(t *testing.T) {
	tests := []struct {
		name     string
		comment  string
		expected []IgnoreComment
	}{
		{
			name:    "single ignore code",
			comment: "gitlab-ci-verify ignore:CODE1",
			expected: []IgnoreComment{
				{
					Comment: "gitlab-ci-verify ignore:CODE1",
					Code:    "CODE1",
				},
			},
		},
		{
			name:    "multiple ignore codes",
			comment: "gitlab-ci-verify ignore:CODE1  ignore:CODE2",
			expected: []IgnoreComment{
				{
					Comment: "gitlab-ci-verify ignore:CODE1  ignore:CODE2",
					Code:    "CODE1",
				},
				{
					Comment: "gitlab-ci-verify ignore:CODE1  ignore:CODE2",
					Code:    "CODE2",
				},
			},
		},
		{
			name:     "no ignore code",
			comment:  "some other comment",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseIgnoreComment(tt.comment)
			if result == nil && tt.expected != nil || result != nil && tt.expected == nil {
				t.Errorf("ParseIgnoreComment(%q) = %v, want %v", tt.comment, result, tt.expected)
			} else if result != nil && tt.expected != nil && !reflect.DeepEqual(tt.expected, result) {
				t.Errorf("ParseIgnoreComment(%q) = %v, want %v", tt.comment, result, tt.expected)
			}
		})
	}
}

func TestParseIgnoreForLine(t *testing.T) {
	tests := []struct {
		name              string
		lineNumberMapping yaml.LineNumberMapping
		lines             []int
		expected          []IgnoreComment
	}{
		{
			name: "single line with ignore comment",
			lineNumberMapping: yaml.LineNumberMapping{
				1: []*yaml.Node{
					{LineComment: "gitlab-ci-verify ignore:CODE1"},
				},
			},
			lines: []int{1},
			expected: []IgnoreComment{
				{Comment: "gitlab-ci-verify ignore:CODE1", Code: "CODE1"},
			},
		},
		{
			name: "multiple lines with ignore comments",
			lineNumberMapping: yaml.LineNumberMapping{
				1: []*yaml.Node{
					{LineComment: "gitlab-ci-verify ignore:CODE1"},
				},
				2: []*yaml.Node{
					{LineComment: "gitlab-ci-verify ignore:CODE2"},
				},
			},
			lines: []int{1, 2},
			expected: []IgnoreComment{
				{Comment: "gitlab-ci-verify ignore:CODE1", Code: "CODE1"},
				{Comment: "gitlab-ci-verify ignore:CODE2", Code: "CODE2"},
			},
		},
		{
			name: "line without ignore comment",
			lineNumberMapping: yaml.LineNumberMapping{
				1: []*yaml.Node{
					{LineComment: "some other comment"},
				},
			},
			lines:    []int{1},
			expected: []IgnoreComment{},
		},
		{
			name: "multiple nodes with ignore comments",
			lineNumberMapping: yaml.LineNumberMapping{
				1: []*yaml.Node{
					{LineComment: "gitlab-ci-verify ignore:CODE1"},
					{LineComment: "gitlab-ci-verify ignore:CODE3"},
				},
				2: []*yaml.Node{
					{LineComment: "gitlab-ci-verify ignore:CODE2"},
				},
			},
			lines: []int{1, 2},
			expected: []IgnoreComment{
				{Comment: "gitlab-ci-verify ignore:CODE1", Code: "CODE1"},
				{Comment: "gitlab-ci-verify ignore:CODE3", Code: "CODE3"},
				{Comment: "gitlab-ci-verify ignore:CODE2", Code: "CODE2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseIgnoreForLine(tt.lineNumberMapping, tt.lines)
			if len(result) != len(tt.expected) {
				t.Errorf("ParseIgnoreForLine() = %v, want %v", result, tt.expected)
			}
			for i := range result {
				if result[i].Comment != tt.expected[i].Comment || result[i].Code != tt.expected[i].Code {
					t.Errorf("ParseIgnoreForLine() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestIgnoreCommentsToCodes(t *testing.T) {
	tests := []struct {
		name     string
		comments []IgnoreComment
		expected []string
	}{
		{
			name: "single ignore code",
			comments: []IgnoreComment{
				{Comment: "gitlab-ci-verify ignore:CODE1", Code: "CODE1"},
			},
			expected: []string{"CODE1"},
		},
		{
			name: "multiple ignore codes",
			comments: []IgnoreComment{
				{Comment: "gitlab-ci-verify ignore:CODE1", Code: "CODE1"},
				{Comment: "gitlab-ci-verify ignore:CODE2", Code: "CODE2"},
			},
			expected: []string{"CODE1", "CODE2"},
		},
		{
			name:     "no ignore codes",
			comments: []IgnoreComment{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IgnoreCommentsToCodes(tt.comments)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("IgnoreCommentsToCodes() = %v, want %v", result, tt.expected)
			}
		})
	}
}
