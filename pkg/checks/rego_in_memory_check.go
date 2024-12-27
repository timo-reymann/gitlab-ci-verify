package checks

import (
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/internal/rego_policies"
)

type InMemoryCheck struct {
	RegoContent string
}

func (im InMemoryCheck) convertToCheckFinding(raw map[string]any) (*CheckFinding, error) {
	return convertToCheckFinding("memory", raw)
}

func (im InMemoryCheck) Run(i *CheckInput) ([]CheckFinding, error) {
	rpm := rego_policies.NewRegoPolicyManager()

	if err := rpm.LoadModuleFromString("check.rego", im.RegoContent); err != nil {
		return nil, fmt.Errorf("failed to load rego bundle %s: %s", "in-memory", err)
	}

	results, err := queryManagerForFindings(rpm, i)
	if err != nil {
		return nil, err
	}

	return parseResults("in-memory", results)
}
