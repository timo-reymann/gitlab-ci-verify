package format_conversion

import (
	"gopkg.in/yaml.v3"
)

// ParseYamlNode from the given raw content
func ParseYamlNode(content []byte) (*yaml.Node, error) {
	var result yaml.Node
	err := yaml.Unmarshal(content, &result)
	return &result, err
}
