package ci_yaml

import (
	formatconversion "github.com/timo-reymann/gitlab-ci-verify/internal/format-conversion"
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
}

// NewCiYamlFile from byte contents
func NewCiYamlFile(content []byte) (*CiYamlFile, error) {
	doc, err := formatconversion.ParseYamlNode(content)
	if err != nil {
		return nil, err
	}

	var parsed map[string]any
	err = doc.Decode(&parsed)
	if err != nil {
		return nil, err
	}

	return &CiYamlFile{
		FileContent:   content,
		ParsedYamlMap: parsed,
		ParsedYamlDoc: doc,
	}, nil
}
