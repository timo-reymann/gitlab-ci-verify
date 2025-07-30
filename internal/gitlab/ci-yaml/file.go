package ci_yaml

import (
	"bytes"
	formatconversion "github.com/timo-reymann/gitlab-ci-verify/v2/internal/format-conversion"
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/gitlab/ci-yaml/includes"
	"github.com/timo-reymann/gitlab-ci-verify/v2/pkg/filtering"
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
	lineNumberMapping yaml.LineNumberMapping
}

func (cyf *CiYamlFile) GetFileLevelIgnores() []filtering.IgnoreComment {
	if len(cyf.ParsedYamlDoc.Content) == 0 || len(cyf.ParsedYamlDoc.Content[0].Content) == 0 {
		return []filtering.IgnoreComment{}
	}
	return filtering.ParseIgnoreComment(cyf.ParsedYamlDoc.Content[0].Content[0].HeadComment)
}

func (cyf *CiYamlFile) GetLineLevelIgnores(line int) []filtering.IgnoreComment {
	return filtering.ParseIgnoreForLine(cyf.lineNumberMapping, []int{line})
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
		lineNumberMapping: lineNummberMapping,
	}, nil
}
