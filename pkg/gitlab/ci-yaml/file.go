package ci_yaml

import (
	"bytes"
	formatconversion "github.com/timo-reymann/gitlab-ci-verify/internal/format-conversion"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/gitlab/ci-yaml/includes"
	"gopkg.in/yaml.v3"
)

// CiYamlFile represents a parsed CI YAML file
type CiYamlFile struct {
	// FileContent contains the raw file content
	FileContent []byte
	// ParsedYamlMap contains the parsed YAML as map structure
	ParsedYamlMap map[string]any
	// ParsedYamlDoc contains the parsed YAML as node structure
	ParsedYamlDoc *yaml.Node
	// Includes contains all includes in the entry file
	Includes []includes.Include
}

// NewCiYamlFile from byte contents
func NewCiYamlFile(content []byte) (*CiYamlFile, error) {
	doc, err := formatconversion.ParseYamlNode(content)
	if err != nil {
		return nil, err
	}

	var parsed map[string]any
	decoder := yaml.NewDecoder(bytes.NewBuffer(content))
	decoder.UniqueKeys(false)

	if err = decoder.Decode(&parsed); err != nil {
		return nil, err
	}

	parsedIncludes, err := includes.ParseIncludes(doc)
	if err != nil {
		return nil, err
	}

	return &CiYamlFile{
		FileContent:   content,
		ParsedYamlMap: parsed,
		ParsedYamlDoc: doc,
		Includes:      parsedIncludes,
	}, nil
}
