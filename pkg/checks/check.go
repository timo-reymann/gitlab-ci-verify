package checks

import (
	format_conversion "github.com/timo-reymann/gitlab-ci-verify/internal/format-conversion"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/cli"
	ci_yaml "github.com/timo-reymann/gitlab-ci-verify/pkg/gitlab/ci-yaml"
	"gopkg.in/yaml.v3"
)

type CiYaml struct {
	FileContent   []byte
	ParsedYamlMap map[string]any
	ParsedYamlDoc *yaml.Node
}

type CheckInput struct {
	// CiYaml contains the yaml configuration in different representations as loaded from the file system
	CiYaml *CiYaml
	// Configuration used to execute the verifier
	Configuration *cli.Configuration
	// LintAPIResult contains the API response from the GitLab CI Lint API. It might be not set when configuration
	// disabled ci validation inside CI
	LintAPIResult *ci_yaml.VerificationResultWithRemoteInfo
	// MergedCiYaml contains the merged yaml when the lint api result is available
	MergedCiYaml *CiYaml
}

func (c *CheckInput) HasLintAPIResult() bool {
	return c.LintAPIResult != nil
}

func (c *CheckInput) CanProvideMergedYaml() bool {
	return c.LintAPIResult != nil
}

// NewCiYaml from byte contents
func NewCiYaml(content []byte) (*CiYaml, error) {
	doc, err := format_conversion.ParseYamlNode(content)
	if err != nil {
		return nil, err
	}

	var parsed map[string]any
	err = doc.Decode(&parsed)
	if err != nil {
		return nil, err
	}

	return &CiYaml{
		FileContent:   content,
		ParsedYamlMap: parsed,
		ParsedYamlDoc: doc,
	}, nil
}

// Check that runs verifications
type Check interface {
	// Run the check
	Run(i *CheckInput) ([]CheckFinding, error)
}
