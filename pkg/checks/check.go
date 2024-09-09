package checks

import (
	format_conversion "github.com/timo-reymann/gitlab-ci-verify/internal/format-conversion"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/cli"
	"gopkg.in/yaml.v3"
)

type CiYaml struct {
	FileContent   []byte
	ParsedYamlMap map[string]any
	ParsedYamlDoc *yaml.Node
}

type CheckInput struct {
	CiYaml        *CiYaml
	Configuration *cli.Configuration
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
