package format_conversion

import (
	"gopkg.in/yaml.v3"
)

// ParseYaml from the given raw content into a generic map
func ParseYaml(content []byte) (map[string]any, error) {
	result := map[string]any{}
	err := yaml.Unmarshal(content, result)
	return result, err
}
