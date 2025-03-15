package format_conversion

import (
	"bytes"
	"gopkg.in/yaml.v3"
)

// ParseYamlNode from the given raw content
func ParseYamlNode(content []byte) (*yaml.Node, error) {
	var result yaml.Node

	decoder := yaml.NewDecoder(bytes.NewBuffer(content))
	decoder.UniqueKeys(false)
	err := decoder.Decode(&result)
	return &result, err
}
