package ci_yaml

import (
	format_conversion "github.com/timo-reymann/gitlab-ci-verify/pkg/format-conversion"
	"os"
)

func LoadRaw(path string) (map[string]any, error) {
	yamlContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return format_conversion.ParseYaml(yamlContent)
}
