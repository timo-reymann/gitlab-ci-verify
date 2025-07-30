package checks

import (
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/v2/internal/rego_policies"
)

type ModuleCheck struct {
	ModulePath string
}

func (m ModuleCheck) Run(i *CheckInput) ([]CheckFinding, error) {
	rpm := rego_policies.NewRegoPolicyManager()

	if err := rpm.LoadModuleFromFile(m.ModulePath); err != nil {
		return nil, fmt.Errorf("failed to load rego module %s: %s", m.ModulePath, err)
	}

	results, err := queryManagerForFindings(rpm, i)
	if err != nil {
		return nil, err
	}

	return parseResults(i, m.ModulePath, results)
}
