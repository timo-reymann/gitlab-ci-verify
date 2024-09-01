package ci_yaml

import (
	"github.com/timo-reymann/gitlab-ci-verify/internal/format-conversion"
	"os"
)

// LoadRaw parses a given Gitlab CI YAML file into a generic map structure
func LoadRaw(path string) (map[string]any, error) {
	yamlContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return format_conversion.ParseYaml(yamlContent)
}
