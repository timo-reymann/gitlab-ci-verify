package checks

import (
	"github.com/google/go-cmp/cmp"
	"github.com/timo-reymann/gitlab-ci-verify/internal/cli"
	"github.com/timo-reymann/gitlab-ci-verify/internal/gitlab/ci-yaml"
	"os"
	"path"
	"testing"
)

func NewCiYamlFromFile(t *testing.T, fileName string) *ci_yaml.CiYamlFile {
	content, err := os.ReadFile(path.Join(fileName))
	if err != nil {
		t.Fatal(err)
	}
	ciYaml, err := ci_yaml.NewCiYamlFile(content)
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

func createCheckInput(t *testing.T, ciYaml *ci_yaml.CiYamlFile, projectRoot, ciFile string) *CheckInput {
	virtualCiYaml, err := ci_yaml.CreateVirtualCiYamlFile(projectRoot, ciFile, ciYaml)
	if err != nil {
		t.Fatal(err)
	}
	return &CheckInput{
		VirtualCiYaml: virtualCiYaml,
		MergedCiYaml:  virtualCiYaml.Combined,
		Configuration: &cli.Configuration{},
	}
}
