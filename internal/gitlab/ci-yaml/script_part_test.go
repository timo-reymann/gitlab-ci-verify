package ci_yaml

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestScriptPart_SplitContentLines(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []ScriptPartLine
	}{
		{
			name:  "Single line",
			input: "This is a single line",
			expected: []ScriptPartLine{
				{LineContent: "This is a single line", LineNumber: 1, Node: nil},
			},
		},
		{
			name:  "Multiple lines",
			input: "Line 1\nLine 2\nLine 3",
			expected: []ScriptPartLine{
				{LineContent: "Line 1", LineNumber: 1, Node: nil},
				{LineContent: "Line 2", LineNumber: 2, Node: nil},
				{LineContent: "Line 3", LineNumber: 3, Node: nil},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scriptPart := ScriptPart{Content: tc.input}
			result := scriptPart.SplitContentLines()

			if !cmp.Equal(result, tc.expected) {
				t.Errorf("Expected:\n%#v\nGot:\n%#v", tc.expected, result)
			}
		})
	}
}

func TestConcat(t *testing.T) {
	testCases := []struct {
		name           string
		parts          []ScriptPart
		expectedLines  []ScriptPartLine
		expectedScript []byte
	}{
		{
			name:           "Empty parts",
			parts:          []ScriptPart{},
			expectedLines:  []ScriptPartLine{},
			expectedScript: []byte{},
		},
		{
			name:  "Single part",
			parts: []ScriptPart{{Content: "Line 1\nLine 2"}},
			expectedLines: []ScriptPartLine{
				{LineContent: "Line 1", LineNumber: 1, Node: nil},
				{LineContent: "Line 2", LineNumber: 2, Node: nil},
			},
			expectedScript: []byte("Line 1\nLine 2\n"),
		},
		{
			name:  "Multiple parts",
			parts: []ScriptPart{{Content: "Part 1"}, {Content: "Part 2\nPart 3"}},
			expectedLines: []ScriptPartLine{
				{LineContent: "Part 1", LineNumber: 1, Node: nil},
				{LineContent: "Part 2", LineNumber: 2, Node: nil},
				{LineContent: "Part 3", LineNumber: 3, Node: nil},
			},
			expectedScript: []byte("Part 1\nPart 2\nPart 3\n"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lines, script := Concat(tc.parts)

			if !reflect.DeepEqual(lines, tc.expectedLines) {
				t.Errorf("Expected lines:\n%#v\nGot:\n%#v", tc.expectedLines, lines)
			}

			if !bytes.Equal(script, tc.expectedScript) {
				t.Errorf("Expected script:\n%s\nGot:\n%s", tc.expectedScript, script)
			}
		})
	}
}
