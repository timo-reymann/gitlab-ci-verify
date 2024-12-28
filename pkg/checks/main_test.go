package checks

import (
	"github.com/google/go-cmp/cmp"
	"os"
	"path"
	"testing"
)

func NewCiYamlFromFile(t *testing.T, fileName string) *CiYaml {
	content, err := os.ReadFile(path.Join(".", fileName))
	if err != nil {
		t.Fatal(err)
	}
	ciYaml, err := NewCiYaml(content)
	if err != nil {
		t.Fatal(err)
	}
	return ciYaml
}

func CheckMustSucceed(findings []CheckFinding, err error) []CheckFinding {
	if err != nil {
		panic(err)
	}

	return findings
}

func VerifyFindings(t *testing.T, expected []CheckFinding, actual []CheckFinding) {
	if !cmp.Equal(expected, actual) {
		t.Fatalf("Expected %v findings, but got %v", expected, actual)
	}
}
