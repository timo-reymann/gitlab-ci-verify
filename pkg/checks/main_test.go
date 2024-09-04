package checks

import (
	"github.com/google/go-cmp/cmp"
	format_conversion "github.com/timo-reymann/gitlab-ci-verify/internal/format-conversion"
	"os"
	"path"
	"testing"
)

func newCiYamlMock(t *testing.T, content []byte) *CiYaml {
	doc, err := format_conversion.ParseYamlNode(content)
	if err != nil {
		t.Fatal(err)
	}

	var parsed map[string]any
	err = doc.Decode(&parsed)
	if err != nil {
		t.Fatal(err)
	}

	return &CiYaml{
		FileContent:   content,
		ParsedYamlMap: parsed,
		ParsedYamlDoc: doc,
	}
}

func newCiYamlFromFile(t *testing.T, fileName string) *CiYaml {
	content, err := os.ReadFile(path.Join(".", fileName))
	if err != nil {
		t.Fatal(err)
	}
	return newCiYamlMock(t, content)
}

func checkMustSucceed(findings []CheckFinding, err error) []CheckFinding {
	if err != nil {
		panic(err)
	}

	return findings
}

func verifyFindings(t *testing.T, expected []CheckFinding, actual []CheckFinding) {
	if !cmp.Equal(expected, actual) {
		t.Fatalf("Expected %v findings, but got %v", expected, actual)
	}
}
