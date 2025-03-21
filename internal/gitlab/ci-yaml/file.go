package ci_yaml

import (
	"bytes"
	formatconversion "github.com/timo-reymann/gitlab-ci-verify/internal/format-conversion"
	"github.com/timo-reymann/gitlab-ci-verify/internal/gitlab/ci-yaml/includes"
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
	// LineNumberMapping contains the line number node mapping for the parsed YAML
	LineNumberMapping *yaml.LineNumberMapping
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
	decoder.WithLineNumberMapping()

	if err = decoder.Decode(&parsed); err != nil {
		return nil, err
	}
	lineNummberMapping := decoder.LineNumberMapping()

	parsedIncludes, err := includes.ParseIncludes(doc)
	if err != nil {
		return nil, err
	}

	return &CiYamlFile{
		FileContent:       content,
		ParsedYamlMap:     parsed,
		ParsedYamlDoc:     doc,
		Includes:          parsedIncludes,
		LineNumberMapping: &lineNummberMapping,
	}, nil
}
