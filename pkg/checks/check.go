package checks

import (
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

// Check that runs verifications
type Check interface {
	// Run the check
	Run(rep *CheckInput) ([]CheckFinding, error)
}
